package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/lanceryou/confd/loader"
)

type Event struct {
	Op   loader.OpType
	Path string
}

type EventHandler func(event Event) error
type FileWatcher struct {
	watcher *fsnotify.Watcher

	watchedMu sync.RWMutex
	watched   map[string]EventHandler

	noExistPath chan string
	eventsDone  chan struct{}
	closeOnce   sync.Once
}

func NewFileWatcher() (fileWatcher *FileWatcher, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	fileWatcher = &FileWatcher{
		watcher:     watcher,
		watched:     make(map[string]EventHandler),
		eventsDone:  make(chan struct{}, 1),
		noExistPath: make(chan string, 1),
	}

	go fileWatcher.runEvents()
	go fileWatcher.runWatchPath()
	return
}

func (w *FileWatcher) runWatchPath() {
	noExistPaths := make(map[string]struct{})
	for {
		select {
		case p := <-w.noExistPath:
			noExistPaths[p] = struct{}{}
		case <-time.After(time.Second):
			var add []string
			for k := range noExistPaths {
				if !isExists(k) {
					continue
				}

				add = append(add, k)
			}

			for _, p := range add {
				if err := w.addWatch(p); err == nil {
					delete(noExistPaths, p)
				}
			}
		}
	}
}

func isExists(path string) bool {
	_, err := os.Stat(path) // os.Stat获取文件信息
	if err != nil && !os.IsExist(err) {
		return false
	}
	return true
}

func (w *FileWatcher) runEvents() {
	defer close(w.eventsDone)

	go func() {
		for err := range w.watcher.Errors {
			log.Println(err)
		}
	}()

	for e := range w.watcher.Events {
		switch {
		case e.Op&fsnotify.Create == fsnotify.Create:
			w.sendEvent(Event{loader.Create, e.Name})
		case e.Op&fsnotify.Write == fsnotify.Write,
			e.Op&fsnotify.Chmod == fsnotify.Chmod:
			w.sendEvent(Event{loader.Update, e.Name})
		case e.Op&fsnotify.Remove == fsnotify.Remove:
			w.sendEvent(Event{loader.Delete, e.Name})
		case e.Op&fsnotify.Rename == fsnotify.Rename:
			// Rename is only issued on the original file path; the new name receives a Create event
			w.sendEvent(Event{loader.Delete, e.Name})
		default:
			panic(fmt.Sprintf("unknown op type %v", e.Op))
		}
	}
}

func (w *FileWatcher) sendEvent(event Event) {
	w.watchedMu.RLock()
	watch, ok := w.watched[event.Path]
	w.watchedMu.RUnlock()

	if !ok {
		d := filepath.Dir(event.Path)
		w.watchedMu.RLock()
		watch, ok = w.watched[d]
		w.watchedMu.RUnlock()
		if !ok {
			return
		}
	}

	watch(event)
}

func (w *FileWatcher) AddWatch(path string, fn EventHandler) (err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	w.watchedMu.RLock()
	_, ok := w.watched[absPath]
	w.watchedMu.RUnlock()
	if ok {
		return
	}

	w.watchedMu.Lock()

	defer w.watchedMu.Unlock()
	w.watched[absPath] = fn

	err = w.addWatch(absPath)
	if err != nil {
		delete(w.watched, absPath)
	}
	return
}

func (w *FileWatcher) addWatch(path string) error {
	if !isExists(path) {
		w.noExistPath <- path
		return nil
	}
	err := w.watcher.Add(path)
	if err != nil && os.IsPermission(err) {
		return fmt.Errorf(fmt.Sprintf("path %v IsPermission", path))
	}

	return err
}

func (w *FileWatcher) RemoveWatch(path string) {
	w.watchedMu.Lock()
	delete(w.watched, path)
	w.watchedMu.Unlock()
}

func (w *FileWatcher) Close() {
	w.closeOnce.Do(func() {
		if w.watcher != nil {
			w.watcher.Close()
			<-w.eventsDone
		}
	})
}
