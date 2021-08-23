package main

import (
	"bitcoin_trade/app/controllers"
	"bitcoin_trade/app/models"
	"bitcoin_trade/config"
	"bitcoin_trade/utils"
	"fmt"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	fmt.Println(models.DbConnection)
	// controllers.StreamIngestionData()
	controllers.StartWebServer()
}
