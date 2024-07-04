package services

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/dto/daemon"
	"go-tha-utxos/config"
	"os/exec"
)

const (
	generateAddressesCmd    = "getnewaddress"
	listAddressGroupingsCmd = "listaddressgroupings"
	listUnspentCmd          = "listunspent"
	createRawTxCmd          = "createrawtransaction"
	signRawTxCmd            = "signrawtransactionwithwallet"
	sendRawTxCmd            = "sendrawtransaction"

	MatureConfirmations = "10"
	MinUnspentAmount    = 0.3
)

type CliCommands struct {
	DaemonCli string
}

func NewCliCommands(cfg *config.Config) (*CliCommands, error) {
	cli := cfg.Commands.DaemonCli
	if cli == "" {
		return nil, fmt.Errorf("config file does not contain command.daemon_cli")
	}

	return &CliCommands{DaemonCli: cli}, nil
}

func (d *CliCommands) GetNewAddresses(count int) (addresses []string, err error) {
	successCalls := 0
	for i := 0; i < count; i++ {
		out, err := exec.Command(d.DaemonCli, generateAddressesCmd).CombinedOutput()
		if err != nil {
			err = fmt.Errorf("command [%s] failed with error: %v. output: %s", generateAddressesCmd, err, string(out))

			return addresses, err
		}

		address := string(out)
		log.Info("generated address: ", address)

		addresses = append(addresses, RemoveLineBreaks(address))

		successCalls++
	}

	return
}

func (d *CliCommands) GetExistingAddresses() (*daemon.ListAddressGroupingsResponse, error) {
	var resp daemon.ListAddressGroupingsResponse
	out, err := exec.Command(d.DaemonCli, listAddressGroupingsCmd).CombinedOutput()
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

	out, err := exec.Command(d.DaemonCli, listUnspentCmd, "1", "9999999", "[]", "false", string(reqAsString)).CombinedOutput()
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when called daemon with error: %v", listUnspentCmd, err)
	}

	//DO NOT COMMIT
	//testing := `[{"txid":"f7e68572588502da54476e67564b88f3b42794c86e4181dff0339be0f38e38c5","vout":0,"address":"1FZ78z3wmdMmHxXqj2Kn67XQyo9mriVi5y","label":"","scriptPubKey":"76a9149fa431f8b7fc6b05c49f300e106cb4f92d66140788ac","amount":344.64709652,"confirmations":5,"spendable":true,"solvable":true,"desc":"pkh([cb3feb02/0h/0h/0h]0200b6748dda7b4660c96d459efe5d03b7a2925b22d56aecbeaf1457715895b06b)#j2akdk4p","parent_descs":[],"safe":true},{"txid":"f7e68572588502da54476e67564b88f3b42794c86e4181dff0339be0f38e38c5","vout":0,"address":"1FZ78z3wmdMmHxXqj2Kn67XQyo9mriVi5y","label":"","scriptPubKey":"76a9149fa431f8b7fc6b05c49f300e106cb4f92d66140788ac","amount":1344.64709652,"confirmations":5,"spendable":true,"solvable":true,"desc":"pkh([cb3feb02/0h/0h/0h]0200b6748dda7b4660c96d459efe5d03b7a2925b22d56aecbeaf1457715895b06b)#j2akdk4p","parent_descs":[],"safe":true}]`
	//out = []byte(testing)

	err = json.Unmarshal(out, &unspent)
	if err != nil {
		return unspent, fmt.Errorf("command [%s] failed when unmarshalling daemon response: %v", listUnspentCmd, err)
	}

	return unspent, nil
}

func (d *CliCommands) CreateRawTransaction(inputs []daemon.Unspent, outputs []daemon.RawTransactionOutput) (rawTx string, err error) {
	unspentToSend := make(map[string]float64, len(inputs))
	for _, input := range inputs {
		unspentToSend[input.Address] = input.Amount
	}

	inputsStr, err := json.Marshal(unspentToSend)
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when marshalling tx inputs: %v", createRawTxCmd, err)
	}

	outputsStr, err := json.Marshal(outputs)
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when marshalling tx outputs: %v", createRawTxCmd, err)
	}

	log.Debugf("calling [%s] with [%s] and [%s] as arguments", createRawTxCmd, string(inputsStr), string(outputsStr))
	out, err := exec.Command(d.DaemonCli, createRawTxCmd, string(inputsStr), string(outputsStr)).CombinedOutput()
	if err != nil {
		return rawTx, fmt.Errorf("command [%s] failed when called daemon with error: %v", listUnspentCmd, err)
	}

	return string(out), nil
}
