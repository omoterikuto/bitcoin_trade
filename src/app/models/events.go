package models

import (
	"encoding/json"
	"fmt"
	"sort"
	"src/config"
	"time"
)

type SignalEvent struct {
	Time        time.Time `json:"time" gorm:"primaryKey"`
	ProductCode string    `json:"product_code" gorm:"type:varchar(191)"`
	Side        string    `json:"side" gorm:"type:varchar(191)"`
	Price       float64   `json:"price" gorm:"type:float"`
	Size        float64   `json:"size"  gorm:"type:float"`
}

type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
}

func NewSignalEvents() *SignalEvents {
	return &SignalEvents{}
}

func GetSignalEventsByCount(loadEvents int) *SignalEvents {
	// cmd := fmt.Sprintf(`SELECT * FROM (
	//     SELECT time, product_code, side, price, size FROM %s WHERE product_code = ? ORDER BY time DESC LIMIT ? )
	//     ORDER BY time ASC;`, tableNameSignalEvents)
	// rows, err := DbConnection.Query(cmd, config.Config.ProductCode, loadEvents)

	eventSlices := []SignalEvent{}
	Db.Where("product_code=?", config.Config.ProductCode).Order("time desc").Limit(loadEvents).Find(&eventSlices)

	sort.Slice(eventSlices, func(i, j int) bool {
		return eventSlices[i].Time.Before(eventSlices[j].Time)
	})

	var signalEvents SignalEvents
	signalEvents.Signals = eventSlices

	return &signalEvents
}

func GetSignalEventsAfterTime(dateTime time.Time) *SignalEvents {
	// cmd := fmt.Sprintf(`SELECT * FROM (
	//             SELECT time, product_code, side, price, size FROM %s
	//             WHERE DATETIME(time) >= DATETIME(?)
	//             ORDER BY time DESC
	//         ) ORDER BY time ASC;`, tableNameSignalEvents)

	// rows, err := DbConnection.Query(cmd, timeTime.Format(time.RFC3339))
	// if err != nil {
	// 	return nil
	// }
	// defer rows.Close()

	// var signalEvents SignalEvents
	// for rows.Next() {
	// 	var signalEvent SignalEvent
	// 	rows.Scan(&signalEvent.Time, &signalEvent.ProductCode, &signalEvent.Side, &signalEvent.Price, &signalEvent.Size)
	// 	signalEvents.Signals = append(signalEvents.Signals, signalEvent)
	// }
	eventSlices := []SignalEvent{}
	Db.Where("time >= ?", dateTime).Order("time desc").Find(&eventSlices)

	sort.Slice(eventSlices, func(i, j int) bool {
		return eventSlices[i].Time.Before(eventSlices[j].Time)
	})

	var signalEvents SignalEvents
	signalEvents.Signals = eventSlices

	return &signalEvents
}

func (s *SignalEvents) CanBuy(time time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return true
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "SELL" && lastSignal.Time.Before(time) {
		return true
	}
	return false
}

func (s *SignalEvents) CanSell(time time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return false
	}

	lastSignal := s.Signals[lenSignals-1]
	if lastSignal.Side == "BUY" && lastSignal.Time.Before(time) {
		return true
	}
	return false
}

func (s *SignalEvents) Buy(ProductCode string, time time.Time, price, size float64, save bool) bool {
	if !s.CanBuy(time) {
		return false
	}
	signalEvent := SignalEvent{
		ProductCode: ProductCode,
		Time:        time,
		Side:        "BUY",
		Price:       price,
		Size:        size,
	}
	if save {
		result := Db.Create(&signalEvent)
		if result.Error != nil {
			fmt.Println(result.Error)
		}
	}
	s.Signals = append(s.Signals, signalEvent)
	return true
}

func (s *SignalEvents) Sell(productCode string, time time.Time, price, size float64, save bool) bool {

	if !s.CanSell(time) {
		return false
	}

	signalEvent := SignalEvent{
		ProductCode: productCode,
		Time:        time,
		Side:        "SELL",
		Price:       price,
		Size:        size,
	}

	if save {
		result := Db.Create(&signalEvent)
		if result.Error != nil {
			fmt.Println(result.Error)
		}
	}

	s.Signals = append(s.Signals, signalEvent)
	return true
}

func (s *SignalEvents) Profit() float64 {
	total := 0.0
	beforeSell := 0.0
	isHolding := false
	for i, signalEvent := range s.Signals {
		if i == 0 && signalEvent.Side == "SELL" {
			continue
		}
		if signalEvent.Side == "BUY" {
			total -= signalEvent.Price * signalEvent.Size
			isHolding = true
		}
		if signalEvent.Side == "SELL" {
			total += signalEvent.Price * signalEvent.Size
			isHolding = false
			beforeSell = total
		}
	}
	if isHolding {
		return beforeSell
	}
	return total
}

func (s SignalEvents) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(&struct {
		Signals []SignalEvent `json:"signals,omitempty"`
		Profit  float64       `json:"profit,omitempty"`
	}{
		Signals: s.Signals,
		Profit:  s.Profit(),
	})
	if err != nil {
		return nil, err
	}
	return value, err
}

func (s *SignalEvents) CollectAfter(time time.Time) *SignalEvents {
	for i, signal := range s.Signals {
		if time.After(signal.Time) {
			continue
		}
		return &SignalEvents{Signals: s.Signals[i:]}
	}
	return nil
}
