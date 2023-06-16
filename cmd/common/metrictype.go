package common

import "net/http"

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

func CheckPostHTTPMethod(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST methods", http.StatusMethodNotAllowed)
		return
	}
}

func CheckGetHTTPMethod(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET methods", http.StatusMethodNotAllowed)
		return
	}
}
