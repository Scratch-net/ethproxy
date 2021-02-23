package service

import (
	"errors"

	"github.com/gorilla/mux"
)

func NewRouter(handler *Handler) (*mux.Router, error) {

	if handler == nil {
		return nil, errors.New("handler is nil")
	}

	router := mux.NewRouter()

	// we accept either "latest" or a 1-20 digit number as a block and a 1-4 digit number or a hash as a tx

	router.HandleFunc("/block/{block:latest|\\d{1,20}}/txs/{tx:\\d{1,4}|0x[0-9a-f]{64}}", Handle(handler.HandleRequest)).
		Methods("GET")

	return router, nil
}
