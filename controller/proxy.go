package controller

import (
	"strconv"

	"github.com/scratch-net/ethproxy/api"
	"github.com/scratch-net/ethproxy/cache"
)

type EthProxyController struct {
	blockCache cache.BlockCache
}

func New(cache cache.BlockCache) (*EthProxyController, error) {

	if cache == nil {
		return nil, api.ErrCacheNotInitialized
	}
	return &EthProxyController{
		blockCache: cache,
	}, nil
}

// FetchTransaction retrieves a block and tries to get tx either by index or by hash
func (e *EthProxyController) FetchTransaction(blockNumber, tx string) (*api.Transaction, error) {
	block, err := e.blockCache.Get(blockNumber)
	if err != nil {
		return nil, err
	}

	// first we try to parse tx as a number
	txNum, err := strconv.Atoi(tx)
	if err != nil {
		// if it's not a number then it's a hash (as guaranteed by router)
		transaction, ok := block.TxByHash[tx]
		if ok {
			return transaction, nil
		}
		return nil, api.ErrTransactionNotFound
	}
	if txNum < len(block.Transactions) {
		return block.Transactions[txNum], nil
	}
	return nil, api.ErrTransactionNotFound
}
