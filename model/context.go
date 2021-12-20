package model

import (
	"context"
	"net/http"

	ua "github.com/mileusna/useragent"
)

type OS struct{}

var osKey OS

func NewContext(r *http.Request) context.Context {
	ctx := r.Context()
	userAgent := r.UserAgent()
	ua := ua.Parse(userAgent)

	return context.WithValue(ctx, osKey, ua)
}
