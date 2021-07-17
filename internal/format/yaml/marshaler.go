package yaml

import (
	"github.com/go-yaml/yaml"
)

type MarshalerYml struct{}

func (m *MarshalerYml) Marshal(src interface{}) ([]byte, error) {
	return yaml.Marshal(src)
}

func (m *MarshalerYml) Unmarshal(data []byte, dst interface{}) error {
	return yaml.Unmarshal(data, dst)
}

func (m *MarshalerYml) String() string {
	return "yaml"
}
