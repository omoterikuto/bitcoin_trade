package main

import (
	"fmt"
	"src/app/controllers"
	"src/app/models"
	"src/config"
	"src/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	// controllers.StreamIngestionData()
	fmt.Println("___________")
	controllers.StartWebServer()

	sqlDb, err := models.Db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDb.Close()
}
