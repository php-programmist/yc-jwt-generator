package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"overseer/api"
)

type errorMessage struct {
	ErrorMessage string `json:"error"`
}

func handleError(w http.ResponseWriter, err error) {
	message := errorMessage{err.Error()}
	jsonData, err2 := json.Marshal(message)
	if err2 != nil {
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	_, err2 = w.Write(jsonData)
	if err != nil {
		return
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var credentials api.Credentials
	w.Header().Set("Content-Type", "application/json")

	err := decodeJSONBody(w, r, &credentials)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var Token, err2 = api.GetIAMToken(credentials)
	if err2 != nil {
		handleError(w, err2)
		return
	}

	data, err3 := json.Marshal(Token)
	if err3 != nil {
		handleError(w, err3)
		return
	}

	_, err4 := w.Write(data)
	if err4 != nil {
		handleError(w, err4)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRequest)

	log.Print("Starting server on :80...")
	err := http.ListenAndServe(":80", mux)
	log.Fatal(err)
}
