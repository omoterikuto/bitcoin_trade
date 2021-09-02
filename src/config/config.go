package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	LogFile          string
	ProductCode      string
	TradeDuration    time.Duration
	Durations        map[string]time.Duration
	BackTest         bool
	UsePercent       float64
	DataLimit        int
	StopLimitPercent float64
	NumRanking       int

	DbName      string
	SQLDriver   string
	DbUser      string
	DbPort      int
	DbPassword  string
	DbContainer string

	WebPort int

	ApiKey    string
	ApiSecret string
}

var Config ConfigList

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
	Config = ConfigList{
		LogFile:          cfg.Section("btc_trade").Key("log_file").String(),
		ProductCode:      cfg.Section("btc_trade").Key("product_code").String(),
		Durations:        durations,
		TradeDuration:    durations[cfg.Section("btc_trade").Key("trade_duration").String()],
		BackTest:         cfg.Section("btc_trade").Key("back_test").MustBool(),
		UsePercent:       cfg.Section("btc_trade").Key("use_percent").MustFloat64(),
		DataLimit:        cfg.Section("btc_trade").Key("data_limit").MustInt(),
		StopLimitPercent: cfg.Section("btc_trade").Key("stop_limit_percent").MustFloat64(),
		NumRanking:       cfg.Section("btc_trade").Key("num_ranking").MustInt(),

		DbName:      cfg.Section("db").Key("name").String(),
		SQLDriver:   cfg.Section("db").Key("driver").String(),
		DbUser:      cfg.Section("db").Key("user").String(),
		DbPort:      cfg.Section("db").Key("port").MustInt(),
		DbPassword:  cfg.Section("db").Key("password").String(),
		DbContainer: cfg.Section("db").Key("container").String(),

		WebPort: cfg.Section("web").Key("port").MustInt(),

		ApiKey:    cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret: cfg.Section("bitflyer").Key("api_secret").String(),
	}
}
