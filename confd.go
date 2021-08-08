package confd

import (
	"fmt"
)

type Options struct {
	schema string
	confer string
	watch  bool
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

func WithWatch(watch bool) OptionFunc {
	return func(o *Options) {
		o.watch = watch
	}
}

// Confd conf manager
type Confd struct {
	opt Options
	Config
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
// loader 加载配置，
// config 只提供读取相关？
func (c *Confd) LoadConfig() (err error) {
	sche, err := ParseSchema(c.opt.schema)
	if err != nil {
		return
	}

	if err = c.loadConfig(sche); err != nil {
		return
	}

	if c.opt.watch {
		c.watchSchema(sche)
	}
	return
}

func (c *Confd) loadConfig(sche *Schema) (err error) {
	loader, err := GetConfigLoader(sche.source)
	if err != nil {
		return
	}

	data, err := loader.Load(sche.key)
	if err != nil {
		return
	}

	marshal, err := GetMarshaler(sche.format)
	if err != nil {
		return
	}

	cfg, err := GetConfig(c.opt.confer)
	if err != nil {
		return
	}

	if err = cfg.ReadIn(marshal, data); err != nil {
		return
	}

	c.Config = cfg
	return
}

// TODO how to watch
func (c *Confd) watchSchema(sche *Schema) {

}
