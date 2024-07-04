package action

import (
	"encoding/json"
	"fmt"
	"go-tha-utxos/app/math"

	//_ "github.com/cosmos/cosmos-sdk/types"
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/dto/daemon"
	"go-tha-utxos/app/dto/response"
	"go-tha-utxos/app/services"
	"go-tha-utxos/config"
	"slices"
)

const (
	minUtxoAmount = 0.1
)

func CreateUtxos(cfg *config.Config, file string, fee float64) error {
	cli, err := services.NewCliCommands(cfg)
	if err != nil {
		log.Fatal(err)
	}

	addresses, err := getAddressesFromFile(file)
	if err != nil {
		return err
	}

	numOfAddresses := len(addresses)
	if numOfAddresses == 0 {
		return fmt.Errorf("no addresses found in file [%s]", file)
	}

	log.Infof("found [%d] addresses in file [%s]. Can try to generate [%d] UTXOs", numOfAddresses, file, numOfAddresses)

	unspent, err := cli.ListUnspent(500)
	if err != nil {
		return fmt.Errorf("finding unspent addresses failed: %s", err)
	}

	if len(unspent) == 0 {
		log.Infof("no usable unspent found")
		return nil
	}

	slices.SortStableFunc(unspent, func(a, b daemon.Unspent) int {
		return int(b.Amount) - int(a.Amount)
	})

	//amountNeeded := minUtxoAmount * float64(numOfAddresses) + fee
	minUtxoDec, _ := math.NewDecFromStr("0.1")
	amountNeeded := minUtxoDec.MulInt64(int64(numOfAddresses))
	feeDec, err := math.NewDecFromStr(services.FloatToString(fee))
	if err != nil {
		return fmt.Errorf("error on fee conversion [%.2f]: %s", fee, err)
	}
	amountNeeded = amountNeeded.Add(feeDec)

	//extract only needed UTXOs
	amountFound := math.ZeroDec()
	var selectedUnspent []daemon.Unspent
	for i := 0; i < len(unspent); i++ {
		if amountFound.GT(amountNeeded) {
			break
		}

		unspendAmount, err := math.NewDecFromStr(services.FloatToString(unspent[i].Amount))
		if err != nil {
			return fmt.Errorf("error on unspend amount conversion [%.2f]: %s", unspendAmount, err)
		}
		amountFound = amountFound.Add(unspendAmount)
		selectedUnspent = append(selectedUnspent, unspent[i])
	}

	log.Infof("amount [%.2f] found in unspent utxos. Needed [%.2f]", amountFound, amountNeeded)

	//subtract the fee and use the rest to create outputs

	amountFound = amountFound.Sub(feeDec)
	//if amountFound < minUtxoAmount {
	if amountFound.LT(minUtxoDec) {
		log.Infof("found amount is lower than minimum [%.2f] needed to create UTXO", amountFound.MustFloat64())
	}

	var outputs []daemon.RawTransactionOutput
	for i := 0; i < numOfAddresses; i++ {
		if amountFound.LT(minUtxoDec) {
			break
		}

		raw := make(daemon.RawTransactionOutput, 1)
		raw[addresses[i]] = minUtxoAmount
		outputs = append(outputs, raw)
		amountFound = amountFound.Sub(minUtxoDec)
	}

	//send the remaining back to the first address
	if amountFound.IsPositive() {
		raw := make(daemon.RawTransactionOutput, 1)
		raw[unspent[0].Address] = amountFound.MustFloat64()
		outputs = append(outputs, raw)
	}

	return nil
}

func getAddressesFromFile(file string) ([]string, error) {
	if !services.FileExists(file) {
		return []string{}, fmt.Errorf("file [%s] does not exist", file)
	}

	var data response.GenerateAddressResponse
	fileData, err := services.ReadFile(file)
	if err != nil {
		return []string{}, err
	}

	if fileData == nil {
		return []string{}, fmt.Errorf("file [%s] is empty", file)
	}

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return []string{}, err
	}

	return data.Addresses, nil
}
