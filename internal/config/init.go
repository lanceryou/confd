package config

import (
	"github.com/lanceryou/confd"
)

func init() {
	confd.RegisterConfig(NewConfig())
}
