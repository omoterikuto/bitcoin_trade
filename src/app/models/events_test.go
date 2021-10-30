package models

import (
	"testing"
	"time"
)

func TestCanBuyAndCanSell(t *testing.T) {
	events := DecodeJSON("testdata/events.json")[0].Events
	signalEvents := SignalEvents{Signals: events}

	t.Run("canBuy", func(t *testing.T) {
		if !signalEvents.CanBuy(time.Now()) {
			t.Error("test failed")
		}
	})
	t.Run("cannotSell", func(t *testing.T) {
		if signalEvents.CanSell(time.Now()) {
			t.Error("test failed")
		}
	})
	signalEvents.Signals = signalEvents.Signals[:3]
	t.Run("cannotBuy", func(t *testing.T) {
		if signalEvents.CanBuy(time.Now()) {
			t.Error("test failed")
		}
	})
	t.Run("canSell", func(t *testing.T) {
		if !signalEvents.CanSell(time.Now()) {
			t.Error("test failed")
		}
	})
}

func TestBuyAndSell(t *testing.T) {
	events := DecodeJSON("testdata/events.json")[0].Events
	signalEvents := SignalEvents{Signals: events}

	t.Run("buy", func(t *testing.T) {
		thisTime := time.Now()
		signalEvents.Buy("BTC_JPY", thisTime, 5000000, 0.1, false)
		// 一番最後のevent情報の時間が一致すれば成功
		if signalEvents.Signals[len(signalEvents.Signals)-1].Time != thisTime {
			t.Error("failed to buy")
		}
	})

	t.Run("sell", func(t *testing.T) {
		thisTime := time.Now()
		signalEvents.Sell("BTC_JPY", thisTime, 5000000, 0.1, false)
		if signalEvents.Signals[len(signalEvents.Signals)-1].Time != thisTime {
			t.Error("failed to sell")
		}
	})
}

func TestProfit(t *testing.T) {
	events := DecodeJSON("testdata/events.json")[0].Events
	signalEvents := SignalEvents{Signals: events}

	if signalEvents.Profit() != 100000 {
		t.Errorf("profit not equal 1000000, but equal %0.0f", signalEvents.Profit())
	}
}
func TestCollectAfter(t *testing.T) {
	events := DecodeJSON("testdata/events.json")[0].Events
	signalEvents := SignalEvents{Signals: events}

	time, _ := time.Parse(time.RFC3339, "2021-01-02T00:00:00+09:00")
	collectEvents := signalEvents.CollectAfter(time)
	if len(collectEvents.Signals) != 3 {
		t.Errorf("failed to collect: signal num is %d", len(collectEvents.Signals))
	}
}
