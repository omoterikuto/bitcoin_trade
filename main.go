package main

import (
	"btc_trade/app/controllers"
	"btc_trade/config"
	"btc_trade/utils"
	"log"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	controllers.StreamIngestionData()
	log.Println(controllers.StartWebServer())
}
