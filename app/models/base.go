package models

import (
	"btc_trade/config"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	tableNameSignalEvents = "signal_events"
)

func GetCandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration)
}

var Db *gorm.DB

func init() {
	fmt.Println("base")
	Db := sqlConnect()

	defer Db.Close()

	// Db.CreateTable(&SignalEvent{})
	fmt.Println("base_end")
	// cmd := fmt.Sprintf(`
	//     CREATE TABLE IF NOT EXISTS %s (
	//         time DATETIME PRIMARY KEY NOT NULL,
	//         product_code STRING,
	//         side STRING,
	//         price FLOAT,
	//         size FLOAT)`, tableNameSignalEvents)
	// DbConnection.Exec(cmd)

	// for _, duration := range config.Config.Durations {
	// 	tableName := GetCandleTableName(config.Config.ProductCode, duration)
	// 	c := fmt.Sprintf(`
	//         CREATE TABLE IF NOT EXISTS %s (
	//         time DATETIME PRIMARY KEY NOT NULL,
	//         open FLOAT,
	//         close FLOAT,
	//         high FLOAT,
	//         low FLOAT,
	// 		volume FLOAT)`, tableName)
	// 	DbConnection.Exec(c)
	// }
}

func sqlConnect() (db *gorm.DB) {
	c := config.Config
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.DbUser, c.DbPassword, c.DbContainer, c.DbPort, c.DbName)
	fmt.Println(dataSourceName)

	count := 0
	db, err := gorm.Open(c.SQLDriver, dataSourceName)

	if err != nil {
		for {
			if err == nil {
				fmt.Println("DB接続成功")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 30 {
				fmt.Println("")
				fmt.Println("DB接続失敗")
				panic(err)
			}
			db, err = gorm.Open(c.SQLDriver, dataSourceName)
		}
	}

	return db
}
