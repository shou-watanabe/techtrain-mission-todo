package middleware

import (
	"net/http"
)

func Layers(handler http.Handler) http.Handler {
	return Recovery(
		Context(
			Access(
				handler,
			),
		),
	)
}
