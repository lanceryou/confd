package file

import (
	"io/ioutil"
	"os"

	"github.com/lanceryou/confd/loader"
)

type FileLoader struct {
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

// TODO 后续实现
func (f *FileLoader) Watch(path string) (ret *loader.WatchResult, err error) {
	return nil, nil
}

// String file loader
func (f *FileLoader) String() string {
	return "file"
}
