package format

import (
	"github.com/lanceryou/confd"
	"github.com/lanceryou/confd/internal/format/json"
	"github.com/lanceryou/confd/internal/format/yaml"
)

func init() {
	confd.RegisterMarshaler(
		&json.MarshalerJson{},
		&yaml.MarshalerYml{},
	)
}
