package config

import (
	"fmt"

	"github.com/lanceryou/confd/format"
)

type ConfigFactory interface {
	New() Config
	String() string
}

// Config config support multi-config format(yml, json,toml ...)
type Config interface {
	Clean() error
	ReadIn(marshaler format.Marshaler, data []byte) error
	Read(val interface{}) error
	ReadSection(key string, val interface{}) error
	HasValue(key string) bool
}

var (
	cpMap = make(map[string]ConfigFactory)
)

// RegisterConfig
func RegisterConfigFactory(c ConfigFactory) {
	cpMap[c.String()] = c
}

// DRegisterConfig deregister config
func DRegisterConfigFactory(name string) {
	delete(cpMap, name)
}

// GetConfig get config by name
func GetConfig(name string) (Config, error) {
	cfg, ok := cpMap[name]
	if !ok {
		return nil, fmt.Errorf("config %v has not register.", name)
	}

	return cfg.New(), nil
}
