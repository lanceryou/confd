package confd

import (
	"fmt"
)

// Marshaler marshaler
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}

var (
	formatMap = make(map[string]Marshaler)
)

// RegisterMarshaler
func RegisterMarshaler(ms ...Marshaler) {
	for _, m := range ms {
		formatMap[m.String()] = m
	}
}

// DRegisterMarshaler deregister marshaler
func DRegisterMarshaler(name string) {
	delete(formatMap, name)
}

// GetMarshaler get marshaler by format
func GetMarshaler(format string) (Marshaler, error) {
	ms, ok := formatMap[format]
	if !ok {
		return nil, fmt.Errorf("marshaler %v has not register.", format)
	}

	return ms, nil
}
