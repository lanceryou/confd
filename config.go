package confd

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
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

type ConfigE struct {
	marshaler Marshaler
	loader    ConfigLoader
	schema    Schema
	config    map[string]interface{}
	watch     bool
}

// LoadConfig load config
func (c *ConfigE) LoadConfig() error {
	data, err := c.loader.Load(c.schema.key)
	if err != nil {
		return err
	}

	err = c.marshaler.Unmarshal(data, &c.config)
	if err != nil {
		return err
	}

	if c.watch {
		// Watch config change
	}

	return nil
}

func (c *ConfigE) Read(val interface{}) error {
	return decode(c.config, val)
}

func (c *ConfigE) ReadSection(key string, val interface{}) error {
	return nil
}

func (c *ConfigE) HasValue(key string) bool {
	return false
}

func decode(input interface{}, output interface{}) error {
	cfg := &mapstructure.DecoderConfig{
		Result: output,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}

	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
