package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	daemon2 "go-tha-utxos/app/dto/daemon"
	daemon "go-tha-utxos/app/dto/daemon/rpc"
	"go-tha-utxos/config"
	"io/ioutil"
	"net/http"
)

type RpcDaemon struct {
	cfg    *config.RpcConnection
	client *http.Client
}

func NewRpcDaemon(rpcConn *config.RpcConnection) (*RpcDaemon, error) {
	if rpcConn == nil {
		return nil, fmt.Errorf("rpc configuration is not valid")
	}

	if err := rpcConn.Validate(); err != nil {
		return nil, err
	}

	return &RpcDaemon{
		client: &http.Client{},
		cfg:    rpcConn,
	}, nil
}

func (bd *RpcDaemon) sendRequest(baseRequest daemon.BaseRequest) ([]byte, error) {
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
func (bd *RpcDaemon) GetNewAddresses(count int) ([]string, error) {
	var addresses []string
	for i := 0; i < count; i++ {
		result, err := bd.sendRequest(daemon.NewBaseRequest("getnewaddress", []interface{}{}))
		if err != nil {
			return nil, err
		}

		var baseResponse daemon.BaseResponse
		if err := json.Unmarshal(result, &baseResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal address: %v", err)
		}

		address := baseResponse.Result
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (bd *RpcDaemon) DumpPrivateKey(address string) (key string, err error) {
	//TODO implement me
	panic("implement me")
}

func (bd *RpcDaemon) GetExistingAddresses() (*daemon2.ListAddressGroupingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (bd *RpcDaemon) ListUnspent(count int) (unspent []daemon2.Unspent, err error) {
	//TODO implement me
	panic("implement me")
}

func (bd *RpcDaemon) CreateRawTransaction(inputs []daemon2.RawTransactionInput, outputs []daemon2.RawTransactionOutput) (rawTx string, err error) {
	//TODO implement me
	panic("implement me")
}

func (bd *RpcDaemon) SignRawTransaction(rawTx string) (signed string, err error) {
	//TODO implement me
	panic("implement me")
}

func (bd *RpcDaemon) SendRawTransaction(hexString string) (txHash string, err error) {
	//TODO implement me
	panic("implement me")
}

func (bd *RpcDaemon) ListUnspentDust(count int) (unspent []daemon2.Unspent, err error) {
	//TODO implement me
	panic("implement me")
}
