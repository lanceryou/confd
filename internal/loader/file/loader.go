package file

import (
	"io/ioutil"
	"os"

	"github.com/lanceryou/confd/loader"
)

type FileLoader struct {
	watcher *FileWatcher
}

// NewFileLoader
func NewFileLoader() *FileLoader {
	watcher, err := NewFileWatcher()
	if err != nil {
		panic(err)
	}

	return &FileLoader{watcher: watcher}
}

// Load load config data
func (f *FileLoader) Load(path string) (data []byte, err error) {
	fs, err := os.Open(path)
	if err != nil {
		return
	}

	defer fs.Close()
	return ioutil.ReadAll(fs)
}

// watch path file change and reload the changed file
func (f *FileLoader) Watch(path string) (ret *loader.WatchResult, err error) {
	retChan := make(chan loader.OpType, 1)
	err = f.watcher.AddWatch(path, func(event Event) error {
		retChan <- event.Op
		return nil
	})
	if err != nil {
		return
	}

	ev := <-retChan
	data, err := f.Load(path)
	if err != nil {
		return
	}
	ret = &loader.WatchResult{Action: ev, Result: data}
	return ret, nil
}

// String file loader
func (f *FileLoader) String() string {
	return "file"
}
