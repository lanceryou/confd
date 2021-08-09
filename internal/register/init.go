package register

import (
	"github.com/lanceryou/confd/config"
	"github.com/lanceryou/confd/format"
	"github.com/lanceryou/confd/internal/config/vp"
	"github.com/lanceryou/confd/internal/format/json"
	"github.com/lanceryou/confd/internal/format/yaml"
	"github.com/lanceryou/confd/internal/loader/file"
	"github.com/lanceryou/confd/loader"
)

func init() {
	config.RegisterConfig(vp.NewConfig())
	format.RegisterMarshaler(
		&json.MarshalerJson{},
		&yaml.MarshalerYml{},
	)
	loader.RegisterConfigLoader(&file.FileLoader{})
}
