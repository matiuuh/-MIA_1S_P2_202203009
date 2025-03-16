package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"proyecto1/Analyzer"
)

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	// Manejar solicitud OPTIONS para CORS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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
