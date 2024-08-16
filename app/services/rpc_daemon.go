package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	daemon "go-tha-utxos/app/dto/daemon/rpc"
	"go-tha-utxos/config"
	"io/ioutil"
	"net/http"
)

type BitcoinDaemon struct {
	cfg    *config.RpcConnection
	client *http.Client
}

func NewBitcoinDaemon(rpcConn *config.RpcConnection) (*BitcoinDaemon, error) {
	if rpcConn == nil {
		return nil, fmt.Errorf("rpc configuration is not valid")
	}

	if err := rpcConn.Validate(); err != nil {
		return nil, err
	}

	return &BitcoinDaemon{
		client: &http.Client{},
		cfg:    rpcConn,
	}, nil
}

func (bd *BitcoinDaemon) sendRequest(baseRequest daemon.BaseRequest) ([]byte, error) {
	requestBody, err := json.Marshal(baseRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", bd.cfg.Host, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	req.SetBasicAuth(bd.cfg.User, bd.cfg.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := bd.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform RPC request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

// GetNewAddresses generates new addresses.
func (bd *BitcoinDaemon) GetNewAddresses(count int) ([]string, error) {
	var addresses []string
	for i := 0; i < count; i++ {
		result, err := bd.sendRequest(daemon.NewBaseRequest("getnewaddress", []interface{}{}))
		if err != nil {
			return nil, err
		}

		var address string
		if err := json.Unmarshal(result, &address); err != nil {
			return nil, fmt.Errorf("failed to unmarshal address: %v", err)
		}

		addresses = append(addresses, address)
	}

	return addresses, nil
}
