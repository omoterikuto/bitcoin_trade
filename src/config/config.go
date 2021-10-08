package config

import (
	"log"
	"time"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	LogFile     string
	ProductCode string
	Durations   map[string]time.Duration

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
	setTimeZone()

	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("Failed to read file: ", err)
	}

	durations := map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
		"1d": time.Hour * 24,
	}

	Config = ConfigList{
		LogFile:     cfg.Section("btc_trade").Key("log_file").String(),
		ProductCode: cfg.Section("btc_trade").Key("product_code").String(),
		Durations:   durations,

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

var Location *time.Location

func setTimeZone() {
	var err error
	Location, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		Location = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
	time.Local = Location
}
