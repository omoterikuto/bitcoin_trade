package models

import (
	"fmt"
	"src/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbCandle struct {
	Time   time.Time `gorm:"primaryKey"`
	Open   float64   `gorm:"float, type:not null"`
	Close  float64   `gorm:"float, type:not null"`
	High   float64   `gorm:"float, type:not null"`
	Low    float64   `gorm:"float, type:not null"`
	Volume float64   `gorm:"float, type:not null"`
}

func GetCandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration)
}

var Db *gorm.DB

func init() {
	fmt.Println("base")
	var err error
	Db, err = sqlConnect()
	if err != nil {
		fmt.Println(err)
	}

	sqlDb, err := Db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDb.Close()

	err = migrate()
	if err != nil {
		fmt.Println(err)
	}
}

func sqlConnect() (sqlDb *gorm.DB, err error) {
	c := config.Config
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.DbUser, c.DbPassword, c.DbContainer, c.DbPort, c.DbName)

	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	count := 0

	if err != nil {
		for {
			if err == nil {
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 10 {
				fmt.Println("DB接続失敗")
				return nil, err
			}
			db, err = gorm.Open(mysql.New(mysql.Config{
				DriverName: "mysql",
				DSN:        dataSourceName,
			}), &gorm.Config{})
		}
	}

	fmt.Println("DB接続成功")

	return db, err
}

func migrate() (err error) {
	Db.AutoMigrate(&SignalEvent{})

	for _, duration := range config.Config.Durations {
		tableName := GetCandleTableName(config.Config.ProductCode, duration)
		if !Db.Migrator().HasTable(tableName) {
			err := Db.Migrator().CreateTable(&DbCandle{})
			if err != nil {
				return err
			}
			err = Db.Migrator().RenameTable(&DbCandle{}, tableName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
