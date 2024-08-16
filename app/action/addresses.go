package action

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/dto/daemon"
	"go-tha-utxos/app/dto/response"
	"go-tha-utxos/app/services"
	"go-tha-utxos/config"
	"time"
)

type AddressesDaemon interface {
	GetNewAddresses(count int) (addresses []string, err error)
	DumpPrivateKey(address string) (key string, err error)
	GetExistingAddresses() (*daemon.ListAddressGroupingsResponse, error)
}

func getAddressesDaemon(cfg *config.Config) (AddressesDaemon, error) {
	d, err := services.NewCliDaemon(cfg)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func CollectAddresses(cfg *config.Config, counter int, file string, ignoreExistingFile bool) error {
	addressesDaemon, err := getAddressesDaemon(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if !ignoreExistingFile && file != "" && services.FileExists(file) {
		log.Fatal(fmt.Errorf("file [%s] already exists", file))
	}

	addresses, err := addressesDaemon.GetExistingAddresses()
	if err != nil {
		log.Fatal(err)
	}

	result := response.GenerateAddressResponse{}
	countedAddresses := 0
	for _, addr := range addresses.Addresses {
		result.Addresses = append(result.Addresses, addr.Address)
		countedAddresses++
		//counter can be 0 if he wants all existing addresses and NO new ones
		if countedAddresses == counter {
			break
		}
	}

	//counter can be 0 if he wants all existing addresses and NO new ones
	diff := counter - countedAddresses
	if diff > 0 {
		newAddr, err := addressesDaemon.GetNewAddresses(diff)
		if err != nil {
			log.Fatal(fmt.Errorf("error getting new addresses: %w", err))
		}

		result.Addresses = append(result.Addresses, newAddr...)
	}

	fileActionResponse(result, file, ignoreExistingFile)

	log.Infof("successfully written addresses to [%s] file", file)

	return nil
}

func GenerateAddresses(cfg *config.Config, counter int, file string, ignoreExistingFile bool) error {
	addressesDaemon, err := getAddressesDaemon(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if !ignoreExistingFile && file != "" && services.FileExists(file) {
		log.Fatal(fmt.Errorf("file [%s] already exists", file))
	}

	var result response.GenerateAddressResponse
	result.Addresses, err = addressesDaemon.GetNewAddresses(counter)
	if err != nil {
		if len(result.Addresses) > 0 {
			log.Warnf("successfully called called node command [%d] times but the process failed before it could create [%d] addresses", len(result.Addresses), counter)
		}

		log.Fatal(err)
	}

	log.Infof("successfully generated %d addresses", len(result.Addresses))

	fileActionResponse(result, file, ignoreExistingFile)

	log.Infof("successfully written addresses to [%s] file", file)

	return nil
}

func ExportAddresses(cfg *config.Config, file string) error {
	addressesDaemon, err := getAddressesDaemon(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if !services.FileExists(file) {
		log.Fatal(fmt.Errorf("file [%s] doesn't exist", file))
	}

	addresses, err := getAddressesFromFile(file)
	if err != nil {
		return err
	}

	numOfAddresses := len(addresses)
	if numOfAddresses == 0 {
		return fmt.Errorf("no addresses found in file [%s]", file)
	}

	var result response.AddressesExportResponse
	for _, addr := range addresses {
		pk, err := addressesDaemon.DumpPrivateKey(addr)
		if err != nil {
			return err
		}

		item := make(response.AddressesExportItem, 1)
		item[addr] = pk
		result = append(result, item)
	}

	log.Infof("successfully exported %d addresses", numOfAddresses)
	now := time.Now().String()
	exportFileName := fmt.Sprintf("export_at_%s_from_file_%s", now, file)
	fileActionResponse(result, exportFileName, false)

	log.Infof("successfully exported addresses to [%s] file", exportFileName)

	return nil
}

func fileActionResponse(result interface{}, file string, ignoreExistingFile bool) {
	marshalled, err := json.Marshal(result)
	if err != nil {
		log.Fatal(fmt.Errorf("result marshalling failed: %v", err))
	}

	if file == "" {
		fmt.Println(string(marshalled))

		return
	}

	err = services.WriteToFileIfNotExists(file, marshalled, ignoreExistingFile)
	if err != nil {
		log.Fatal(err)
	}
}
