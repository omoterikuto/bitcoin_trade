package models

import (
	"testing"
)

func TestTableName(t *testing.T) {
	candle := DecodeJSON("testdata/candle.json")[0].Candle
	if candle.TableName() != "BTC_JPY_1m0s" {
		t.Errorf("expected table name %s does not match BTC_JPY_1m0s", candle.TableName())
	}
}
