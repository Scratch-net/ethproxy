package cache

import (
	"time"
)

func CalcTTL(blockNumber, latestBlockNumber uint64) time.Duration {

	// case when we didn't fetch the latest block in time
	if blockNumber > latestBlockNumber {
		return DefaultTTL
	}

	distance := latestBlockNumber - blockNumber

	// Last 20 blocks have small TTL due to possible reorg
	if distance <= 20 {
		// the further the block from the last one, the more TTL it has.
		return DefaultTTL
	}

	// between 20 and 1000 blocks we set TTL depending on distance. The further the block the longer its TTL
	if distance <= 1000 {
		return DefaultTTL * time.Duration(distance)
	}

	// blocks that are safe to cache get 10 years of TTL
	return time.Hour * 24 * 365 * 10
}
