package vp

import (
	"github.com/spf13/viper"

	"github.com/lanceryou/confd/config"
	"github.com/lanceryou/confd/format"
)

func NewConfig() config.ConfigFactory {
	return &viperConfig{vp: viper.New()}
}

type viperConfig struct {
	vp *viper.Viper
}

func (v *viperConfig) New() config.Config {
	return &viperConfig{vp: viper.New()}
}

func (v *viperConfig) Clean() error {
	v.vp = viper.New()
	return nil
}

func (v *viperConfig) ReadIn(marshaler format.Marshaler, data []byte) error {
	cfg := make(map[string]interface{})
	if err := marshaler.Unmarshal(data, &cfg); err != nil {
		return err
	}

	return v.vp.MergeConfigMap(cfg)
}

func (v *viperConfig) Read(val interface{}) error {
	return v.vp.Unmarshal(val)
}

func (v *viperConfig) ReadSection(key string, val interface{}) error {
	return v.vp.UnmarshalKey(key, val)
}

func (v *viperConfig) HasValue(key string) bool {
	return v.vp.IsSet(key)
}

func (v *viperConfig) String() string {
	return "default"
}
