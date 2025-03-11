package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"proyecto1/Analyzer"
)

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Entrada string `json:"entrada"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultado := Analyzer.Analyze(body.Entrada)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"resultado": resultado,
	})
}

func main() {
	http.HandleFunc("/analyze", analyzeHandler)
	fmt.Println("Servidor escuchando en :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}