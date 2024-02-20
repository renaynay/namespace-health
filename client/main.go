package main

import (
	"encoding/json"
	"net/http"
)

type Data struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		data := Data{Message: "Hello, World!"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	// Serve HTML files in the current directory
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
