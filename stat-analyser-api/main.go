package main

import (
	"net/http"

	"github.com/daku-souru/stat-analyser-api/controller"
)

func main() {
	http.HandleFunc("/stat-analyser-api", controller.ProcessVerdictRequest)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
