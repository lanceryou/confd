package confd

import (
	"fmt"
)

// Config config support multi-config format(yml, json,toml ...)
type Config interface {
	Clean() error
	ReadIn(marshaler Marshaler, data []byte) error
	Read(val interface{}) error
	ReadSection(key string, val interface{}) error
	HasValue(key string) bool
	String() string
}

var (
	cpMap = make(map[string]Config)
)

// RegisterConfig
func RegisterConfig(c Config) {
	cpMap[c.String()] = c
}

// DRegisterConfig deregister config
func DRegisterConfig(name string) {
	delete(cpMap, name)
}

// GetConfig get config by name
func GetConfig(name string) (Config, error) {
	cfg, ok := cpMap[name]
	if !ok {
		return nil, fmt.Errorf("config %v has not register.", name)
	}

	return cfg, nil
}
