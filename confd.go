package confd

type Options struct {
	loader string
	format string
	confer string
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
func (c *Confd) LoadConfig(schema string) (err error) {
	sche, err := ParseSchema(schema)
	if err != nil {
		return
	}

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
