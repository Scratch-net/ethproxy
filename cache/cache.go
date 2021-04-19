package cache

import (
	"math"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/karlseguin/ccache/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"

	"github.com/scratch-net/ethproxy/api"
	"github.com/scratch-net/ethproxy/fetcher"
)

const (
	DefaultTTL  = time.Second * 5
	LatestBlock = "latest"
)

type BlockCache interface {
	Get(key string) (*api.Block, error)
}

type BlockCacheStorage struct {
	cache             *ccache.Cache
	fetcher           fetcher.BlockFetcher
	singleFlightGroup singleflight.Group
}

func NewBlockCacheStorage(fetcher fetcher.BlockFetcher, cacheSize int64) (BlockCache, error) {

	if fetcher == nil {
		return nil, api.ErrFetcherNotInitialized
	}
	return &BlockCacheStorage{
		cache:             ccache.New(ccache.Configure().MaxSize(cacheSize)),
		fetcher:           fetcher,
		singleFlightGroup: singleflight.Group{},
	}, nil
}

func (c *BlockCacheStorage) Get(key string) (*api.Block, error) {

	latest, err := c.internalGet(LatestBlock, DefaultTTL)

	if err != nil {
		return nil, err
	}

	if key == LatestBlock {
		return latest, nil
	}

	blockNumber, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		log.Errorf("error while parsing block number: %+v", err)
		return nil, api.ErrInternalServerError
	}

	if blockNumber > math.MaxInt64 {
		return nil, api.ErrBlockNumberTooHigh
	}

	latestBlockNumber, err := hexutil.DecodeUint64(latest.Number)
	if err != nil {
		log.Errorf("error while parsing latest block number: %+v", err)
		return nil, api.ErrInternalServerError
	}

	ttl := CalcTTL(blockNumber, latestBlockNumber)

	return c.internalGet(hexutil.EncodeUint64(blockNumber), ttl)
}

func (c *BlockCacheStorage) internalGet(key string, ttl time.Duration) (*api.Block, error) {
	if c.fetcher == nil {
		return nil, api.ErrFetcherNotSet
	}

	// fetchSingleFunc limits requests to one per block in case multiple clients request the same block simultaneously
	fetchSingleFunc := func() (interface{}, error) {
		res, err, _ := c.singleFlightGroup.Do(key, func() (interface{}, error) {
			block, err := c.fetcher.FetchBlock(key)
			if err == nil {
				c.cache.Set(key, block, ttl)
			}
			return block, err
		})
		return res, err
	}

	item := c.cache.Get(key)

	// no item in cache, fetch it
	if item == nil {
		newItem, err := fetchSingleFunc()
		if err == nil {
			log.Debugf("cache miss: %s", key)
			return newItem.(*api.Block), nil
		}
		return nil, err
	}

	// item exists but expired. Fetch a new one in background, then return the current one
	if item.Expired() {
		go func() {
			_, err := fetchSingleFunc()
			if err != nil {
				c.cache.Delete(key) // unable to fetch fresh item, so delete old one to not to confuse users
				log.Errorf("error while fetching block %s: %+v", key, err)
			}
		}()
	}

	log.Debugf("cache hit: %s", key)
	return item.Value().(*api.Block), nil
}
