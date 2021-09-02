package conf

import (
	"github.com/caarlos0/env"
	"sync"
	"time"
)

type EnvironmentConfig struct {
	LogLevel             string        `env:"ABESH_LOG_LEVEL" envDefault:"debug"`
	CMDLogEnabled        bool          `env:"ABESH_CMD_LOG_ENABLED" envDefault:"false"`
	EventBufferSize      int           `env:"ABESH_EVENT_BUFFER_SIZE" envDefault:"100"`
	GlobalRequestTimeout time.Duration `env:"ABESH_GLOBAL_REQUEST_TIMEOUT" envDefault:"50ms"`
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
