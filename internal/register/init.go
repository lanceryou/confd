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
	config.RegisterConfigFactory(vp.NewConfig())
	format.RegisterMarshaler(
		&json.MarshalerJson{},
		&yaml.MarshalerYaml{},
		&yaml.MarshalerYml{},
	)
	loader.RegisterConfigLoader(file.NewFileLoader())
}
