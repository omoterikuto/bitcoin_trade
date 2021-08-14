package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

//APIInfo bitflyer API
type APIInfo struct {
	APISecret string
	APIKey    string
	LogFile   string
}

//Config APIInfo
var Config APIInfo

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = APIInfo{
		APIKey:    cfg.Section("bitflyer").Key("api_key").String(),
		APISecret: cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:   cfg.Section("fintech").Key("log_file").String(),
	}
}
