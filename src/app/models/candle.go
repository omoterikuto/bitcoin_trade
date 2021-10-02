package models

import (
	"errors"
	"fmt"
	"sort"
	"src/bitflyer"
	"src/config"
	"time"

	"gorm.io/gorm"
)

type Candle struct {
	ProductCode string        `json:"product_code" gorm:"-"`
	Duration    time.Duration `json:"duration" gorm:"-"`
	Time        time.Time     `json:"time" gorm:"primaryKey; datetime"`
	Open        float64       `json:"open" gorm:"not null"`
	Close       float64       `json:"close" gorm:"not null"`
	High        float64       `json:"high" gorm:"not null"`
	Low         float64       `json:"low" gorm:"not null"`
	Volume      float64       `json:"volume" gorm:"not null"`
}

func NewCandle(productCode string, duration time.Duration, dateTime time.Time, open, close, high, low, volume float64) *Candle {
	return &Candle{
		productCode,
		duration,
		dateTime,
		open,
		close,
		high,
		low,
		volume,
	}
}

func (c *Candle) TableName() string {
	return GetCandleTableName(c.ProductCode, c.Duration)
}

func GetCandle(productCode string, duration time.Duration, dateTime time.Time) *Candle {
	tableName := GetCandleTableName(productCode, duration)

	candle := Candle{}

	result := Db.Table(tableName).Where("time = ?", dateTime).First(&candle)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		fmt.Println(result.Error)
		return nil
	}

	candle.Time = dateTime
	candle.Duration = duration

	return &candle
}

func CreateCandleWithDuration(ticker bitflyer.Ticker, productCode string, duration time.Duration) bool {
	currentCandle := GetCandle(productCode, duration, ticker.TruncateDateTime(duration))
	price := ticker.GetMidPrice()
	tableName := GetCandleTableName(config.Config.ProductCode, duration)

	if currentCandle == nil {
		candle := Candle{}
		candle.Time = ticker.TruncateDateTime(duration)
		candle.Open = price
		candle.Close = price
		candle.High = price
		candle.Low = price
		candle.Volume = ticker.Volume

		Db.Table(tableName).Create(&candle)
		return true
	}

	if currentCandle.High <= price {
		currentCandle.High = price
	} else if currentCandle.Low >= price {
		currentCandle.Low = price
	}
	currentCandle.Volume += ticker.Volume
	currentCandle.Close = price

	Db.Table(tableName).Save(&currentCandle)
	return false
}

func GetAllCandle(productCode string, duration time.Duration, limit int) (dfCandle *DataFrameCandle, err error) {
	tableName := GetCandleTableName(productCode, duration)
	candles := []Candle{}

	result := Db.Table(tableName).Order("time desc").Limit(limit).Find(&candles)

	if result.Error != nil {
		return nil, result.Error
	}

	sort.Slice(candles, func(i, j int) bool {
		return candles[i].Time.Before(candles[j].Time)
	})

	dfCandle = &DataFrameCandle{}
	dfCandle.ProductCode = productCode
	dfCandle.Duration = duration
	dfCandle.Candles = candles

	return dfCandle, nil
}
