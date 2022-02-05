package httputils

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	res, err := marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

func RespondErrorJSON(w http.ResponseWriter, status int, err error) {
	he := HTTPError{Message: err.Error()}
	RespondJSON(w, status, he)
}

func marshal(payload interface{}) ([]byte, error) {
	if payload == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(payload)
}
