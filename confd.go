package confd

import (
	"fmt"

	"github.com/lanceryou/confd/config"
	"github.com/lanceryou/confd/format"
	_ "github.com/lanceryou/confd/internal/register"
	"github.com/lanceryou/confd/loader"
)

type Options struct {
	schema string
	confer string
}

func (o *Options) apply() {
	if o.confer == "" {
		o.confer = "default"
	}

	if o.schema == "" {
		panic(fmt.Errorf("schema must be valid"))
	}
}

type OptionFunc func(*Options)

func WithSchema(schema string) OptionFunc {
	return func(o *Options) {
		o.schema = schema
	}
}

func WithConfer(confer string) OptionFunc {
	return func(o *Options) {
		o.confer = confer
	}
}

// Confd conf manager
type Confd struct {
	opt  Options
	sche *Schema
	loader.ConfigLoader
	format.Marshaler
	config.Config
}

func NewConfd(opts ...OptionFunc) *Confd {
	var opt Options
	for _, o := range opts {
		o(&opt)
	}

	opt.apply()
	return &Confd{opt: opt}
}

// LoadConfig load config
// schema {source}:{format}:{key}
// source file, etcd,...
// format ymal, json,ini...
// key can read config information, eg ./xx.yml...
func (c *Confd) LoadConfig() (err error) {
	sche, err := ParseSchema(c.opt.schema)
	if err != nil {
		return
	}

	return c.loadConfig(sche)
}

func (c *Confd) loadConfig(sche *Schema) (err error) {
	ld, err := loader.GetConfigLoader(sche.source)
	if err != nil {
		return
	}

	data, err := ld.Load(sche.key)
	if err != nil {
		return
	}

	marshal, err := format.GetMarshaler(sche.format)
	if err != nil {
		return
	}

	cfg, err := config.GetConfig(c.opt.confer)
	if err != nil {
		return
	}

	if err = cfg.ReadIn(marshal, data); err != nil {
		return
	}

	c.Config = cfg
	c.ConfigLoader = ld
	c.Marshaler = marshal
	return
}

// watch config change.
func (c *Confd) WatchConfig() (err error) {
	for {
		ret, err := c.ConfigLoader.Watch(c.sche.key)
		if err != nil {
			return err
		}

		switch ret.Action {
		case "create", "update":
			if err = c.Config.Clean(); err != nil {
				return err
			}
			err = c.Config.ReadIn(c.Marshaler, ret.Result)
		case "delete":
			err = c.Config.Clean()
		}

		if err != nil {
			return err
		}
	}
}
