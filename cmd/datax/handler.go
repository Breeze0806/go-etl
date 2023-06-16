package main

import (
	"net/http"

	"github.com/Breeze0806/go-etl/datax"
)

type handler struct {
	engine *datax.Engine
}

func newHandler(engine *datax.Engine) *handler {
	return &handler{
		engine: engine,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if h.engine.Metrics().JSON() == nil {
		return
	}
	w.Write([]byte(h.engine.Metrics().JSON().String()))
}
