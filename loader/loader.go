package loader

import (
	"fmt"
)

// the watcher. Actions can be create, update, delete
type WatchResult struct {
	Action string
	Result []byte
}

//  ConfigLoader config load interface
type ConfigLoader interface {
	Load(path string) ([]byte, error)
	Watch(path string) (*WatchResult, error)
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
func GetConfigLoader(name string) (ConfigLoader, error) {
	loader, ok := cfgMap[name]
	if !ok {
		return nil, fmt.Errorf("loader %v has not register.", name)
	}

	return loader, nil
}
