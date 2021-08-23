package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

//APIInfo bitflyer API
type APIInfo struct {
	APISecret     string
	APIKey        string
	LogFile       string
	ProductCode   string
	TradeDuration time.Duration
	Durations     map[string]time.Duration
	DbName        string
	SQLDriver     string
	Port          int
}

//Config APIInfo
var Config APIInfo

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	durations := map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}

	Config = APIInfo{
		APIKey:        cfg.Section("bitflyer").Key("api_key").String(),
		APISecret:     cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:       cfg.Section("fintech").Key("log_file").String(),
		ProductCode:   cfg.Section("fintech").Key("product_code").String(),
		Durations:     durations,
		TradeDuration: durations[cfg.Section("gotrading").Key("trade_duration").String()],
		DbName:        cfg.Section("db").Key("name").String(),
		SQLDriver:     cfg.Section("db").Key("driver").String(),
		Port:          cfg.Section("web").Key("port").MustInt(),
	}
}
