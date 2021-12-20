package middleware

import (
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

func Context(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := model.NewContext(r)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
