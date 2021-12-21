package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func Access(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		end := time.Now()
		sub := end.Sub(start)
		latency := int64(sub / time.Millisecond)
		access := model.NewAccess(r, start, latency)

		json, err := json.Marshal(access)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(json))
	}
	return http.HandlerFunc(fn)
}
