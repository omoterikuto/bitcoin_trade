package main

import (
	"src/app/controllers"
	"src/config"
	"src/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	controllers.StreamIngestionData()
	controllers.StartWebServer()
}
