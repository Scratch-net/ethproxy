package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/scratch-net/ethproxy/cache"
)

func TestTTL(t *testing.T) {

	type TestCase struct {
		blockNumber, latestBlockNumber uint64
		ttl                            time.Duration
	}

	testCases := []*TestCase{
		{
			blockNumber:       100,
			latestBlockNumber: 101,
			ttl:               cache.DefaultTTL,
		},
		{
			blockNumber:       101,
			latestBlockNumber: 100,
			ttl:               cache.DefaultTTL,
		},
		{
			blockNumber:       80,
			latestBlockNumber: 100,
			ttl:               cache.DefaultTTL,
		},

		{
			blockNumber:       79,
			latestBlockNumber: 100,
			ttl:               cache.DefaultTTL * 21,
		},
		{
			blockNumber:       1000,
			latestBlockNumber: 2001,
			ttl:               time.Hour * 24 * 365 * 10,
		},
	}

	for _, tc := range testCases {
		ttl := cache.CalcTTL(tc.blockNumber, tc.latestBlockNumber)
		require.Equal(t, tc.ttl, ttl)
	}

}
