package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorParametr struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-Type", "application/json")
	respBody := errorParametr{Error: msg}
	dat, err := json.Marshal(respBody)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("{\"error\": \"Something went wrong\"}"))
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	dat, err := json.Marshal(payload)

	if err != nil {
		respondWithError(w, code, "Something went wrong")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}