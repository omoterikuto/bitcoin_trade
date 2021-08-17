package main

import (
	"fintech/bitflyer"
	"fintech/config"
	"fintech/utils"
	"fmt"
	"time"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	apiClient := bitflyer.New(config.Config.APIKey, config.Config.APISecret)

	tickerChannel := make(chan bitflyer.Ticker)
	go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerChannel)
	for ticker := range tickerChannel {
		fmt.Println(ticker)
		fmt.Println(ticker.GetMidPrice())
		// fmt.Println(ticker.TruncateDateTime(time.Second))
		// fmt.Println(ticker.TruncateDateTime(time.Minute))
		fmt.Println(ticker.TruncateDateTime(time.Hour))
	}
}
