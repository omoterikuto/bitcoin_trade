package main

import (
	"fmt"
	"src/app/controllers"
	"src/config"
	"src/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	// controllers.StreamIngestionData()
	fmt.Println("___________")
	controllers.StartWebServer()
}
