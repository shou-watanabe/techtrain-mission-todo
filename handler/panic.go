package handler

import (
	"net/http"
)

type PanicHandler struct{}

// NewPanicHandler returns TODOHandler based http.Handler.
func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}

func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ctx := model.NewContext(r)

	panic(http.StatusInternalServerError)
}
