package confd

//  ConfigLoader config load interface
type ConfigLoader interface {
	// schema {file}:{format}:{path}
	Load(path string) ([]byte, error)
	Watch()
	Format() string
	String() string
}

var (
	cfgMap = make(map[string]ConfigLoader)
)

// RegisterConfig register config
func RegisterConfigLoader(c ConfigLoader) {
	cfgMap[c.String()] = c
}

// DRegisterConfig deregister Config
func DRegisterConfigLoader(name string) {
	delete(cfgMap, name)
}

// GetConfig get config
func GetConfigLoader(name string) ConfigLoader {
	return cfgMap[name]
}
