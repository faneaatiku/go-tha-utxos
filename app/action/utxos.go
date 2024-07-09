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
	var selectedUnspent []daemon.RawTransactionInput
	for i := 0; i < len(unspent); i++ {
		if amountFound.GT(amountNeeded) {
			break
		}

		unspendAmount, err := math.NewDecFromStr(services.FloatToString(unspent[i].Amount))
		if err != nil {
			return fmt.Errorf("error on unspend amount conversion [%.2f]: %s", unspendAmount, err)
		}
		amountFound = amountFound.Add(unspendAmount)
		selectedUnspent = append(selectedUnspent, daemon.RawTransactionInput{
			Txid: unspent[i].Txid,
			Vout: unspent[i].Vout,
		})
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
		//risky code that iterates through a map assuming there only one item in it map[address]amount
		//replace that amount to contain 0.1 sent initially + remaining amount
		raw := make(daemon.RawTransactionOutput, 1)
		for address, _ := range outputs[0] {
			//add 0.1 that was already in this output to the remaining amount
			remainingAmountFloat := amountFound.Add(minUtxoDec).MustFloat64()

			//create the RawTransactionOutput map
			raw[address] = services.ToFixedFloat(remainingAmountFloat, 8)

			//no need to continue; just needed to add remaining amount to one of the outputs
			break
		}

		//replace first item with this map
		outputs[0] = raw
	}

	sendTransaction(cli, selectedUnspent, outputs)

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

func ConsolidateUtxos(cfg *config.Config, fee float64, minUtxos int) error {
	cli, err := services.NewCliCommands(cfg)
	if err != nil {
		log.Fatal(err)
	}

	feeDec, err := math.NewDecFromStr(services.FloatToString(fee))
	if err != nil {
		return fmt.Errorf("error on fee conversion [%.2f]: %s", fee, err)
	}

	unspent, err := cli.ListUnspentDust(500)
	if err != nil {
		return fmt.Errorf("finding unspent utxos failed: %s", err)
	}

	numOfUnspent := len(unspent)
	if numOfUnspent == 0 {
		log.Infof("no usable unspent found")
		return nil
	}

	totalUnspent := math.ZeroDec()
	var selectedUnspent []daemon.RawTransactionInput
	for _, u := range unspent {
		unspendAmount, err := math.NewDecFromStr(services.FloatToString(u.Amount))
		if err != nil {
			return fmt.Errorf("error on unspend amount conversion [%.2f]: %s", u.Amount, err)
		}

		totalUnspent = totalUnspent.Add(unspendAmount)
		selectedUnspent = append(selectedUnspent, daemon.RawTransactionInput{
			Txid: u.Txid,
			Vout: u.Vout,
		})
	}

	remainingUnspent := totalUnspent.Sub(feeDec)
	minUtxoDec, _ := math.NewDecFromStr("0.1")
	//(minUtxos + number of already existing UTXOS) * 0.1
	minAmtNeeded := math.NewDec(int64(minUtxos)).Add(math.NewDec(int64(numOfUnspent))).Mul(minUtxoDec)
	if remainingUnspent.LT(minAmtNeeded) {
		log.Infof("nothing to do. found %s dust and needed %s", remainingUnspent.String(), minAmtNeeded.String())

		return nil
	}

	raw := make(daemon.RawTransactionOutput, 1)
	remainingAmountFloat := remainingUnspent.MustFloat64()
	raw[unspent[0].Address] = services.ToFixedFloat(remainingAmountFloat, 8)
	var outputs []daemon.RawTransactionOutput
	outputs = append(outputs, raw)

	sendTransaction(cli, selectedUnspent, outputs)

	return nil
}

func sendTransaction(cli *services.CliCommands, inputs []daemon.RawTransactionInput, outputs []daemon.RawTransactionOutput) {
	if cli == nil {
		log.Fatal("invalid cli provided")

		return
	}

	rawTx, err := cli.CreateRawTransaction(inputs, outputs)
	if err != nil {
		log.Fatal(fmt.Errorf("error on create raw transaction: %s", err))
	}

	if rawTx == "" {
		log.Fatal(fmt.Errorf("create raw transaction returned empty string"))
	}

	signedHex, err := cli.SignRawTransaction(rawTx)
	if err != nil {
		log.Fatal(fmt.Errorf("error on sign raw transaction: %s", err))
	}

	txHash, err := cli.SendRawTransaction(signedHex)
	if err != nil {
		log.Fatal(fmt.Errorf("error on send raw transaction: %s", err))
	}

	fmt.Println("Transaction sent successfully:")
	fmt.Println(txHash)
	fmt.Println("Finished.")
}
