package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	daemon2 "go-tha-utxos/app/dto/daemon"
	daemon "go-tha-utxos/app/dto/daemon/rpc"
	"go-tha-utxos/config"
	"io/ioutil"
	"net/http"
	"strings"
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
		result, err := bd.sendRequest(daemon.NewBaseRequest(generateAddressesCmd, []interface{}{}))
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
	result, err := bd.sendRequest(daemon.NewBaseRequest(dumpPrivKeyCmd, []interface{}{address}))
	if err != nil {
		return "", err
	}

	var baseResponse daemon.BaseResponse
	if err := json.Unmarshal(result, &baseResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal address: %v", err)
	}

	if baseResponse.Error != nil {
		return "", fmt.Errorf("failed to dump private key: %v", baseResponse.Error)
	}

	return baseResponse.Result, nil
}

func (bd *RpcDaemon) GetExistingAddresses() (*daemon2.ListAddressGroupingsResponse, error) {
	result, err := bd.sendRequest(daemon.NewBaseRequest(listAddressGroupingsCmd, []interface{}{}))
	if err != nil {
		return nil, fmt.Errorf("command [%s] failed with error: %v", listAddressGroupingsCmd, err)
	}

	var baseResponse daemon.ListAddressGroupingsResponse
	err = json.Unmarshal(result, &baseResponse)
	if err != nil {
		return nil, fmt.Errorf("command [%s] failed with error: %v", listAddressGroupingsCmd, err)
	}

	resp := &daemon2.ListAddressGroupingsResponse{
		Addresses: baseResponse.Result,
	}

	return resp, nil
}

func (bd *RpcDaemon) ListUnspent(count int) (unspent []daemon2.Unspent, err error) {
	req := daemon2.ListUnspentRequest{
		MaximumCount:  count,
		MinimumAmount: MinUnspentAmount,
	}

	reqAsString, err := json.Marshal(req)
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when marshalling request for daemon: %v", listUnspentCmd, err)
	}

	result, err := bd.sendRequest(daemon.NewBaseRequest(listUnspentCmd, []interface{}{1, 9999999, "[]", "false", string(reqAsString)}))
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when called daemon with error: %v", listUnspentCmd, err)
	}

	var baseResponse daemon.ListUnspentResponse
	err = json.Unmarshal(result, &baseResponse)
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when unmarshalling daemon response: %v", listUnspentCmd, err)
	}

	return baseResponse.Result, nil
}

func (bd *RpcDaemon) CreateRawTransaction(inputs []daemon2.RawTransactionInput, outputs []daemon2.RawTransactionOutput) (rawTx string, err error) {
	inputsStr, err := json.Marshal(inputs)
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when marshalling tx inputs: %v", createRawTxCmd, err)
	}

	outputsStr, err := json.Marshal(outputs)
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when marshalling tx outputs: %v", createRawTxCmd, err)
	}

	log.Debugf("calling [%s] with inputs: %s and outputs: %s", createRawTxCmd, string(inputsStr), string(outputsStr))
	out, err := bd.sendRequest(daemon.NewBaseRequest(createRawTxCmd, []interface{}{string(inputsStr), string(outputsStr)}))
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when called daemon with error: %v", createRawTxCmd, err)
	}
	out = out[:len(out)-1]

	return strings.TrimSpace(string(out)), nil
}

func (bd *RpcDaemon) SignRawTransaction(rawTx string) (signed string, err error) {
	log.Debugf("calling [%s] with [%s] as argument", signRawTxCmd, rawTx)
	out, err := bd.sendRequest(daemon.NewBaseRequest(signRawTxCmd, []interface{}{rawTx}))
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when called daemon with error: %v", signRawTxCmd, err)
	}

	var response daemon2.SignRawTransactionResponse
	err = json.Unmarshal(out, &response)
	if err != nil {
		return "", fmt.Errorf("command [%s] failed when unmarshalling daemon response: %v", signRawTxCmd, err)
	}

	if !response.Complete {
		return "", fmt.Errorf("command [%s] returned NOT complete transaction: %s", signRawTxCmd, string(out))
	}

	return strings.TrimSpace(response.Hex), nil
}

func (bd *RpcDaemon) SendRawTransaction(hexString string) (txHash string, err error) {
	log.Debugf("calling [%s] with [%s] as argument", sendRawTxCmd, hexString)

	out, err := bd.sendRequest(daemon.NewBaseRequest(sendRawTxCmd, []interface{}{hexString}))
	if err != nil {
		return "", fmt.Errorf("command [%s] failed when called daemon with error: %v", sendRawTxCmd, err)
	}

	return strings.TrimSpace(string(out)), nil
}

func (bd *RpcDaemon) ListUnspentDust(count int) (unspent []daemon2.Unspent, err error) {
	maxDust := maxUnspentDustAmount
	req := daemon2.ListUnspentRequest{
		MaximumCount:  count,
		MinimumAmount: minUnspentDustAmount,
		MaximumAmount: &maxDust,
	}

	reqAsString, err := json.Marshal(req)
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when marshalling request for daemon: %v", listUnspentCmd, err)
	}

	result, err := bd.sendRequest(daemon.NewBaseRequest(listUnspentCmd, []interface{}{1, 9999999, "[]", "false", string(reqAsString)}))
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when called daemon with error: %v", listUnspentCmd, err)
	}

	var baseResponse daemon.ListUnspentResponse
	err = json.Unmarshal(result, &baseResponse)
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when unmarshalling daemon response: %v", listUnspentCmd, err)
	}

	return baseResponse.Result, nil
}
