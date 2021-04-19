package service

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/scratch-net/ethproxy/api"
	"github.com/scratch-net/ethproxy/controller"
)

type Handler interface {
	HandleTxRequest(request *http.Request) ([]byte, error)
	HandleBlockRequest(request *http.Request) ([]byte, error)
}

type ProxyHandler struct {
	app controller.ProxyController
}

func NewHandler(ctrl controller.ProxyController) (Handler, error) {
	if ctrl == nil {
		return nil, api.ErrControllerNotSet
	}

	return &ProxyHandler{
		app: ctrl,
	}, nil
}

func (h *ProxyHandler) HandleTxRequest(request *http.Request) ([]byte, error) {
	vars := mux.Vars(request)

	blockNum := vars["block"]
	txIndex := vars["tx"]

	return h.app.FetchTransaction(blockNum, txIndex)
}

func (h *ProxyHandler) HandleBlockRequest(request *http.Request) ([]byte, error) {
	vars := mux.Vars(request)

	blockNum := vars["block"]

	return h.app.FetchBlock(blockNum)
}
