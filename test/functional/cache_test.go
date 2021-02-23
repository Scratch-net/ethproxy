package functional

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/scratch-net/ethproxy/api"
	"github.com/scratch-net/ethproxy/cache"
	mocks "github.com/scratch-net/ethproxy/fetcher/mock"
)

func TestCache_FetchNonLatestBlock_Should_Call_Twice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	fetcher := mocks.NewMockBlockFetcher(mockCtrl)
	call1 := fetcher.EXPECT().FetchBlock("latest")
	call1.DoAndReturn(func(num string) (*api.Block, error) {
		return &api.Block{Number: "0x123"}, nil
	}).Times(1)

	call2 := fetcher.EXPECT().FetchBlock("0x7b")
	call2.DoAndReturn(func(num string) (*api.Block, error) {
		return &api.Block{Number: "0x7b"}, nil
	}).Times(1)
	call2.After(call1)

	c, err := cache.NewBlockCacheStorage(fetcher, 100)
	require.NoError(t, err)
	block, err := c.Get("123")
	require.NoError(t, err)
	require.NotNil(t, block)
}

func TestCache_TooHigh_ShouldError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	fetcher := mocks.NewMockBlockFetcher(mockCtrl)
	call1 := fetcher.EXPECT().FetchBlock("latest")
	call1.DoAndReturn(func(num string) (*api.Block, error) {
		return &api.Block{Number: "0x123"}, nil
	}).Times(1)

	c, err := cache.NewBlockCacheStorage(fetcher, 100)
	require.NoError(t, err)
	block, err := c.Get("9223372036854775808")
	require.Error(t, api.ErrBlockNumberTooHigh)
	require.Nil(t, block)
}
