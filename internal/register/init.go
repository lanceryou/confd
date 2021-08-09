package register

import (
	"github.com/lanceryou/confd"
	"github.com/lanceryou/confd/internal/config/vp"
	"github.com/lanceryou/confd/internal/format/json"
	"github.com/lanceryou/confd/internal/format/yaml"
	"github.com/lanceryou/confd/internal/loader/file"
)

func init() {
	confd.RegisterConfig(vp.NewConfig())
	confd.RegisterMarshaler(
		&json.MarshalerJson{},
		&yaml.MarshalerYml{},
	)
	confd.RegisterConfigLoader(&file.FileLoader{})
}
