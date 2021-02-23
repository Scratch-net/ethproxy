package integration

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/scratch-net/ethproxy/cache"
	"github.com/scratch-net/ethproxy/controller"
	"github.com/scratch-net/ethproxy/fetcher"
	"github.com/scratch-net/ethproxy/service"
)

var expectedTx = []byte("{\"blockHash\":\"0x4b06575b8e0b1905ad10b38ed05c4c9e2e1b76c0eba61f8f71d0ec6a58961ef0\",\"blockNumber\":\"0xb56828\",\"from\":\"0x445051967075cf24fc22aea0dbcaf308400b2d18\",\"gas\":\"0x5208\",\"gasPrice\":\"0x6760765800\",\"hash\":\"0x3a6efcd15a7ca459df934893277f57a4b33e95fe1407773eb4358961a9f060fc\",\"input\":\"0x\",\"nonce\":\"0x0\",\"to\":\"0x6407a44223471195211761610d3b2069f71d57f6\",\"transactionIndex\":\"0x0\",\"value\":\"0x6f984114c61c3d\",\"v\":\"0x26\",\"r\":\"0x85c814eb197005168e3d0f199d7e2f031821646c2f195ba0f400487e759c388c\",\"s\":\"0x250f16c64591462268254436fad4f63e5b3d7ab75b28064712fd23acd4761da0\"}")

func TestService(t *testing.T) {

	cli := fetcher.NewCloudflareEthClient(time.Second * 5)

	cacheStorage, err := cache.NewBlockCacheStorage(cli, 100)
	require.NoError(t, err)

	ctrl, err := controller.New(cacheStorage)
	require.NoError(t, err)
	handler, err := service.NewHandler(ctrl)
	require.NoError(t, err)

	router, err := service.NewRouter(handler)
	require.NoError(t, err)
	server := httptest.NewServer(router)
	defer server.Close()

	path := "/block/11888680/txs/0"
	res, err := http.Get(server.URL + path)
	require.NoError(t, err)
	require.NotNil(t, res)
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotNil(t, body)
	require.Equal(t, expectedTx, body)
}
