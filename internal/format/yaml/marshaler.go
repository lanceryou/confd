package yaml

import (
	"github.com/go-yaml/yaml"
)

type MarshalerYaml struct{}

func (m *MarshalerYaml) Marshal(src interface{}) ([]byte, error) {
	return yaml.Marshal(src)
}

func (m *MarshalerYaml) Unmarshal(data []byte, dst interface{}) error {
	return yaml.Unmarshal(data, dst)
}

func (m *MarshalerYaml) String() string {
	return "yaml"
}

type MarshalerYml struct {
	MarshalerYaml
}

func (m *MarshalerYml) String() string {
	return "yml"
}
