package service

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/scratch-net/ethproxy/api"
	"github.com/scratch-net/ethproxy/controller"
)

type Handler struct {
	app *controller.EthProxyController
}

func NewHandler(ctrl *controller.EthProxyController) (*Handler, error) {
	if ctrl == nil {
		return nil, api.ErrControllerNotSet
	}

	return &Handler{
		app: ctrl,
	}, nil
}

func (h *Handler) HandleRequest(request *http.Request) (*api.Transaction, error) {
	vars := mux.Vars(request)

	blockNum := vars["block"]
	txIndex := vars["tx"]

	return h.app.FetchTransaction(blockNum, txIndex)
}
