package main

import (
	"net/http"

	"github.com/roaris/go-sns-api/httputils"
)

type AppHandler struct {
	h func(http.ResponseWriter, *http.Request) (int, interface{}, error)
}

func (a AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, payload, err := a.h(w, r)
	if err != nil {
		httputils.RespondErrorJSON(w, status, err)
		return
	}
	httputils.RespondJSON(w, status, payload)
}
