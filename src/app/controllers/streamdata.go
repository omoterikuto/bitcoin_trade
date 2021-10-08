package controllers

import (
	"fmt"
	"src/app/models"
	"src/bitflyer"
	"src/config"
)

func StreamIngestionData() {
	c := config.Config
	s := models.TradeSetting{}
	models.Db.First(&s)
	fmt.Println(s)

	ai := NewAI(c.ProductCode, c.Durations[s.TradeDuration], s.DataLimit, s.UseRate, s.StopLimitRate, s.BackTest)

	var tickerChannl = make(chan bitflyer.Ticker)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerChannl)
	go func() {
		for ticker := range tickerChannl {
			for _, duration := range c.Durations {
				isCreated := models.CreateCandleWithDuration(ticker, ticker.ProductCode, duration)
				if isCreated && duration == c.Durations[s.TradeDuration] {
					ai.Trade()
				}
			}
		}
	}()
}
