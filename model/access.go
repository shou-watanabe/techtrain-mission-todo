package model

import (
	"net/http"
	"time"
)

type Access struct {
	Timestamp time.Time
	Latency   int64
	Path      string
	OS        string
}

func NewAccess(r *http.Request, timestamp time.Time, latency int64) *Access {
	return &Access{
		Timestamp: timestamp,
		Latency:   latency,
		Path:      r.URL.Path,
		OS:        r.Context().Value(OsKey).(string),
	}
}
