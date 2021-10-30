package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type TestEvent struct {
	Events []SignalEvent `json:"events"`
	Candle Candle        `json:"candle"`
}

func DecodeJSON(src string) (testEvents []TestEvent) {
	// JSON読み込み
	data, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	if err := json.Unmarshal(data, &testEvents); err != nil {
		log.Fatal(err)
	}
	return
}
