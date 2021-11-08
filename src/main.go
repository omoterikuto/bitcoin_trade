package main

import (
	"src/app/controllers"
	"src/config"
	"src/utils"

	"google.golang.org/appengine"
)

func main() {
	if !appengine.IsAppEngine() {
		utils.LoggingSettings(config.Config.LogFile)
	}
	controllers.StreamIngestionData()
	controllers.StartWebServer()
}
