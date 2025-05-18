package infra

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type RootConfig struct {
	LogDest   string `json:"log_dest"`
	LogLevel  string `json:"log_level"`
	SQLiteDSN string `json:"sqlite_dsn"`

	SMTPPort       int      `json:"smtp_port"`
	Pipeline       string   `json:"pipeline"`
	AllowedSenders []string `json:"allowed_senders"`
	AllowedHosts   []string `json:"allowed_hosts"`
}

func NewRootConfig() (cfg RootConfig, err error) {
	if len(os.Args) > 1 {
		configFile := os.Args[1]

		cfg, err = parseConfig(configFile)
		if err != nil {
			err = fmt.Errorf("failed to load config file from '%s' (error: %v)", configFile, err)
		}
		return
	}
	return cfg, nil
}

func parseConfig(configFile string) (cfg RootConfig, err error) {
	file, err := os.Open(configFile)
	if err != nil {
		return
	}

	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonBytes, &cfg)
	return
}
