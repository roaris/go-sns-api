package middlewares

import (
	"net/http"

	"github.com/roaris/go-sns-api/httputils"
)

type handler func(w http.ResponseWriter, r *http.Request) (int, interface{}, error)

// レスポンスを作成するミドルウェア
func ResponseMiddleware(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, payload, err := h(w, r)
		if err != nil {
			httputils.RespondErrorJSON(w, status, err)
			return
		}
		httputils.RespondJSON(w, status, payload)
	}
}
