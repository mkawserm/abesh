package conf

import (
	"github.com/caarlos0/env"
	"sync"
)

type EnvironmentConfig struct {
	LogLevel        string `env:"ABESH_LOG_LEVEL" envDefault:"debug"`
}

var instantiated *EnvironmentConfig
var once sync.Once

// EnvironmentConfigIns environment config
func EnvironmentConfigIns() *EnvironmentConfig {
	once.Do(func() {
		instantiated = &EnvironmentConfig{}
		if err := env.Parse(instantiated); err != nil {
			panic(err.Error())
		}
	})
	return instantiated
}
