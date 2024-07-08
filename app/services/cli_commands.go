package services

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/dto/daemon"
	"go-tha-utxos/config"
	"os/exec"
	"strings"
)

const (
	generateAddressesCmd    = "getnewaddress"
	listAddressGroupingsCmd = "listaddressgroupings"
	listUnspentCmd          = "listunspent"
	createRawTxCmd          = "createrawtransaction"
	signRawTxCmd            = "signrawtransactionwithwallet"
	sendRawTxCmd            = "sendrawtransaction"
	dumpPrivKeyCmd          = "dumpprivkey"

	MatureConfirmations = "10"
	MinUnspentAmount    = 0.3
)

type CliCommands struct {
	DaemonCli string
	DataDir   string
}

func NewCliCommands(cfg *config.Config) (*CliCommands, error) {
	cli := cfg.Commands.DaemonCli
	if cli == "" {
		return nil, fmt.Errorf("config file does not contain command.daemon_cli")
	}

	if cfg.Commands.DataDir == "" {
		return nil, fmt.Errorf("config file does not contain command.data_dir")
	}

	return &CliCommands{DaemonCli: cli, DataDir: cfg.Commands.DataDir}, nil
}

func (d *CliCommands) GetNewAddresses(count int) (addresses []string, err error) {
	successCalls := 0
	for i := 0; i < count; i++ {
		out, err := exec.Command(d.DaemonCli, d.getDataDir(), generateAddressesCmd).CombinedOutput()
		if err != nil {
			err = fmt.Errorf("command [%s] failed with error: %v. output: %s", generateAddressesCmd, err, string(out))

			return addresses, err
		}

		address := strings.TrimSpace(string(out))
		log.Info("generated address: ", address)

		addresses = append(addresses, RemoveLineBreaks(address))

		successCalls++
	}

	return
}

func (d *CliCommands) DumpPrivateKey(address string) (key string, err error) {
	out, err := exec.Command(d.DaemonCli, d.getDataDir(), dumpPrivKeyCmd, address).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("command [%s] failed with error: %v. output: %s", dumpPrivKeyCmd, err, string(out))

		return key, err
	}

	key = strings.TrimSpace(string(out))
	log.Infof("dumped key for address: [%s]", address)

	return
}

func (d *CliCommands) GetExistingAddresses() (*daemon.ListAddressGroupingsResponse, error) {
	var resp daemon.ListAddressGroupingsResponse
	out, err := exec.Command(d.DaemonCli, d.getDataDir(), listAddressGroupingsCmd).CombinedOutput()
	if err != nil {
		return &resp, fmt.Errorf("command [%s] failed with error: %v", listAddressGroupingsCmd, err)
	}

	err = json.Unmarshal(out, &resp)
	if err != nil {
		return &resp, fmt.Errorf("command [%s] failed with error: %v", listAddressGroupingsCmd, err)
	}

	return &resp, nil
}

func (d *CliCommands) ListUnspent(count int) (unspent []daemon.Unspent, err error) {
	req := daemon.ListUnspentRequest{
		MaximumCount:  count,
		MinimumAmount: MinUnspentAmount,
	}

	reqAsString, err := json.Marshal(req)
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when marshalling request for daemon: %v", listUnspentCmd, err)
	}

	out, err := exec.Command(d.DaemonCli, d.getDataDir(), listUnspentCmd, "1", "9999999", "[]", "false", string(reqAsString)).CombinedOutput()
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when called daemon with error: %v", listUnspentCmd, err)
	}

	err = json.Unmarshal(out, &unspent)
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when unmarshalling daemon response: %v", listUnspentCmd, err)
	}

	return unspent, nil
}

func (d *CliCommands) CreateRawTransaction(inputs []daemon.RawTransactionInput, outputs []daemon.RawTransactionOutput) (rawTx string, err error) {
	inputsStr, err := json.Marshal(inputs)
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when marshalling tx inputs: %v", createRawTxCmd, err)
	}

	outputsStr, err := json.Marshal(outputs)
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when marshalling tx outputs: %v", createRawTxCmd, err)
	}

	log.Debugf("calling [%s] with inputs: %s and outputs: %s", createRawTxCmd, string(inputsStr), string(outputsStr))
	out, err := exec.Command(d.DaemonCli, d.getDataDir(), createRawTxCmd, string(inputsStr), string(outputsStr)).CombinedOutput()
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when called daemon with error: %v", createRawTxCmd, err)
	}
	out = out[:len(out)-1]

	return strings.TrimSpace(string(out)), nil
}

func (d *CliCommands) SignRawTransaction(rawTx string) (signed string, err error) {

	log.Debugf("calling [%s] with [%s] as argument", signRawTxCmd, rawTx)
	out, err := exec.Command(d.DaemonCli, d.getDataDir(), signRawTxCmd, rawTx).CombinedOutput()
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when called daemon with error: %v", signRawTxCmd, err)
	}

	var response daemon.SignRawTransactionResponse
	err = json.Unmarshal(out, &response)
	if err != nil {
		return "", fmt.Errorf("command [%s] failed when unmarshalling daemon response: %v", signRawTxCmd, err)
	}

	if !response.Complete {
		return "", fmt.Errorf("command [%s] returned NOT complete transaction: %s", signRawTxCmd, string(out))
	}

	return strings.TrimSpace(response.Hex), nil
}

func (d *CliCommands) SendRawTransaction(hexString string) (txHash string, err error) {
	log.Debugf("calling [%s] with [%s] as argument", sendRawTxCmd, hexString)

	out, err := exec.Command(d.DaemonCli, d.getDataDir(), sendRawTxCmd, hexString).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command [%s] failed when called daemon with error: %v", sendRawTxCmd, err)
	}

	return strings.TrimSpace(string(out)), nil
}

func (d *CliCommands) getDataDir() string {
	if d.DataDir == "" {
		return ""
	}

	return fmt.Sprintf("-datadir=%s", d.DataDir)
}
