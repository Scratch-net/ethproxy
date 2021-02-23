package service

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/scratch-net/ethproxy/api"
)

type HTTPFunc func(req *http.Request) (*api.Transaction, error)

func Handle(f HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx, err := f(r)

		if err != nil {
			if httpErr, ok := err.(*api.HTTPError); ok {
				w.WriteHeader(httpErr.Code)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			if _, err = io.WriteString(w, err.Error()); err != nil {
				log.Errorf("[Error] can't write response: %+v", err)
				return
			}
			return
		}

		data, err := json.Marshal(tx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("[Error] HTTP 5xx, can't marshal response: %+v", err)

			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(data); err != nil {
			log.Errorf("[Error] can't write response: %+v", err)
			return
		}
	}
}
