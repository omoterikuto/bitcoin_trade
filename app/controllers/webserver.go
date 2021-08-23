package controllers

import (
	"bitcoin_trade/config"
	"fmt"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("app/views/google.html"))

func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "google.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartWebServer() error {
	http.HandleFunc("/chart/", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
