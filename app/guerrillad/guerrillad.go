package guerrillad

import (
	"fmt"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/backends"
)

func NewGuerillaDaemon(logger infra.RootLogger, rootCfg infra.RootConfig) (guerrilla.Daemon, error) {
	d := guerrilla.Daemon{
		Logger: logger,
		Config: toAppConfig(rootCfg),
	}
	return d, nil
}

func toAppConfig(rootCfg infra.RootConfig) *guerrilla.AppConfig {
	return &guerrilla.AppConfig{
		LogLevel: rootCfg.LogLevel,
		Servers: []guerrilla.ServerConfig{
			{
				ListenInterface: fmt.Sprintf("127.0.0.1:%d", rootCfg.SMTPPort),
				IsEnabled:       true,
			},
		},
		BackendConfig: backends.BackendConfig{
			"save_process": rootCfg.Pipeline,
		},
		AllowedHosts: rootCfg.AllowedHosts,
	}
}
