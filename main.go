package main

import (
	"btc_trade/app/controllers"
	"btc_trade/config"
	"btc_trade/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	controllers.StreamIngestionData()
	controllers.StartWebServer()
}
