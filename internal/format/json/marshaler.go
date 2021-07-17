package json

import (
	"encoding/json"
)

type MarshalerJson struct{}

func (m *MarshalerJson) Marshal(src interface{}) ([]byte, error) {
	return json.Marshal(src)
}

func (m *MarshalerJson) Unmarshal(data []byte, dst interface{}) error {
	return json.Unmarshal(data, dst)
}

func (m *MarshalerJson) String() string {
	return "json"
}
