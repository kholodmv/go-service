package main

import (
	"net/http"
)

func main() {
	store := NewMetricsStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", store.PostHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
