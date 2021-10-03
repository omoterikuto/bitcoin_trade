package models

import (
	"errors"
	"fmt"
	"log"
	"src/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TradeSetting struct {
	ID            int     `gorm:"primaryKey"`
	TradeDuration string  `gorm:"default:1m"`
	BackTest      bool    `gorm:"default:true"`
	UseRate       float64 `gorm:"default:0.9"`
	DataLimit     int     `gorm:"default:365"`
	StopLimitRate float64 `gorm:"default:0.1"`
	NumRanking    int     `gorm:"default:3"`
}

func GetCandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration)
}

var Db *gorm.DB

func init() {
	var err error
	Db, err = sqlConnect()
	if err != nil {
		log.Fatal("sql connect error", err)
	} else {
		log.Println("success to connect to database!")
	}

	err = migrate()
	if err != nil {
		log.Println("migrate error", err)
	}
}

func sqlConnect() (sqlDb *gorm.DB, err error) {
	c := config.Config
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s", c.DbUser, c.DbPassword, c.DbContainer, c.DbPort, c.DbName, "Asia%2FTokyo")

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      dataSourceName,
		DisableDatetimePrecision: true,
	}), &gorm.Config{})

	count := 0

	if err != nil {
		for {
			if err == nil {
				break
			}
			log.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 20 {
				return nil, err
			}
			db, err = gorm.Open(mysql.New(mysql.Config{
				DriverName: "mysql",
				DSN:        dataSourceName,
			}), &gorm.Config{})
		}
	}

	return db, err
}

func migrate() (err error) {
	Db.AutoMigrate(&SignalEvent{})
	Db.AutoMigrate(&TradeSetting{})

	result := Db.Take(&TradeSetting{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		Db.Create(&TradeSetting{})
	}

	for _, duration := range config.Config.Durations {
		tableName := GetCandleTableName(config.Config.ProductCode, duration)
		if !Db.Migrator().HasTable(tableName) {
			err := Db.Migrator().CreateTable(&Candle{})
			if err != nil {
				return err
			}
			err = Db.Migrator().RenameTable(&Candle{}, tableName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
