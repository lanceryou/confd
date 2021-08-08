package loader

import (
	"github.com/lanceryou/confd"
	"github.com/lanceryou/confd/internal/loader/file"
)

func init() {
	confd.RegisterConfigLoader(&file.FileLoader{})
}
