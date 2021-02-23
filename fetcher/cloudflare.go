package fetcher

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/scratch-net/ethproxy/api"
)

// BlockFetcher is an interface common for all fetcher implementations
type BlockFetcher interface {
	FetchBlock(key string) (*api.Block, error)
}

const (
	cloudflareHost = "https://cloudflare-eth.com"
	contentType    = "application/json"
)

type CloudflareEthClient struct {
	cli *http.Client
}

func NewCloudflareEthClient(timeout time.Duration) BlockFetcher {

	return &CloudflareEthClient{
		cli: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				IdleConnTimeout:     time.Minute * 5,
				MaxIdleConnsPerHost: 10,
			},
		},
	}
}

func (c *CloudflareEthClient) FetchBlock(number string) (*api.Block, error) {
	log.Debugf("fetching %s..", number)

	requestID := randomID() // as required by JSON RPC

	req := &api.CloudflareRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{number, true},
		ID:      requestID,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Errorf("could not marshal request: %+v", err)
		return nil, api.ErrInternalServerError
	}

	resp, err := c.cli.Post(cloudflareHost, contentType, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Errorf("error while making request: %+v", err)
		return nil, api.ErrBadGateway
	}

	var respBody []byte
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("could not read gateway response: %+v", err)
		return nil, api.ErrBadGateway
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("bad gateway response: [%d] %s %+v", resp.StatusCode, respBody, err)
		return nil, api.ErrBadGateway
	}

	var response *api.Response
	if err = json.Unmarshal(respBody, &response); err != nil {
		log.Errorf("unable to unmarshal response %+v", err)
		return nil, api.ErrBadGateway
	}

	// case when the gateway returns different response from what we've asked
	if response.ID != requestID {
		log.Errorf("Gateway returned mismatched response id: got %d want %d", response.ID, requestID)
		return nil, api.ErrBadGateway
	}

	if response.Block == nil {
		return nil, api.ErrNoSuchBlock
	}

	block := response.Block

	block.TxByHash = make(map[string]*api.Transaction)
	for _, tx := range block.Transactions {
		block.TxByHash[tx.Hash] = tx
	}

	return block, nil
}

// randomId generates random 16 bit id for Json-RPC using go crypto/rand package
func randomID() int {
	b := make([]byte, 2)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return int(binary.BigEndian.Uint16(b))
}
