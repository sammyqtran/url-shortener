package main

import (
	"encoding/json"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func main() {

	http.HandleFunc("/ping", handler)
	http.ListenAndServe(":8080", nil)
}
