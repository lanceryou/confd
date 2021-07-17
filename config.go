package confd

// Config config support multi-config format(yml, json,toml ...)
type Config interface {
	Set(data []byte)

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

// GetConfig get config by format
func GetConfig(format string, data []byte) Config {
	cfg, ok := cpMap[format]
	if !ok {
		return nil
	}

	cfg.Set(data)
	return cfg
}
