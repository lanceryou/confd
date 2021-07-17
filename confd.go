package confd

import (
	"fmt"
)

type Options struct {
	loader string
	format string
}

func (o *Options) apply() {
	if o.loader == "" {
		o.loader = "file"
	}

	if o.format == "" {
		o.format = "yaml"
	}
}

type OptionFunc func(*Options)

func WithLoader(loader string) OptionFunc {
	return func(o *Options) {
		o.loader = loader
	}
}

func WithFormat(format string) OptionFunc {
	return func(o *Options) {
		o.format = format
	}
}

// Confd conf manager
type Confd struct {
	opt Options
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

// 解析schema
func (c *Confd) LoadConfig(schema string) (Config, error) {
	loader := GetConfigLoader(c.opt.loader)
	if loader == nil {
		return nil, fmt.Errorf("loader:%v is not support", c.opt.loader)
	}

	data, err := loader.Load(schema)
	if err != nil {
		return nil, fmt.Errorf("load path %v fail:%v", schema, err)
	}

	return GetConfig(loader.Format(), data), nil
}
