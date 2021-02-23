package api

import (
	"errors"
	"net/http"
)

var (
	ErrInternalServerError   = NewHTTPError(http.StatusInternalServerError, "internal server error")
	ErrNoSuchBlock           = NewHTTPError(http.StatusNotFound, "block not found")
	ErrBlockNumberTooHigh    = NewHTTPError(http.StatusBadRequest, "block number too high")
	ErrBadGateway            = NewHTTPError(http.StatusBadGateway, "bad gateway")
	ErrTransactionNotFound   = NewHTTPError(http.StatusNotFound, "transaction not found")
	ErrFetcherNotSet         = errors.New("fetcher is not set")
	ErrCacheNotInitialized   = errors.New("cache is not initialised")
	ErrFetcherNotInitialized = errors.New("fetcher is not initialised")
	ErrControllerNotSet      = errors.New("controller is not set")
)
