package service

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/scratch-net/ethproxy/api"
)

type HTTPFunc func(req *http.Request) ([]byte, error)

func Handle(f HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := f(r)

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
		_, err = w.Write(res)
		if err != nil {
			log.Errorf("[Error] can't write response: %+v", err)
		}
	}
}
