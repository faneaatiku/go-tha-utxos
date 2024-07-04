package action

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/dto/response"
	"go-tha-utxos/app/services"
	"go-tha-utxos/config"
)

func CollectAddresses(cfg *config.Config, counter int, file string, ignoreExistingFile bool) error {
	cli, err := services.NewCliCommands(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if !ignoreExistingFile && file != "" && services.FileExists(file) {
		log.Fatal(fmt.Errorf("file [%s] already exists", file))
	}

	addresses, err := cli.GetExistingAddresses()
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
		newAddr, err := cli.GetNewAddresses(diff)
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
	cli, err := services.NewCliCommands(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if !ignoreExistingFile && file != "" && services.FileExists(file) {
		log.Fatal(fmt.Errorf("file [%s] already exists", file))
	}

	var result response.GenerateAddressResponse
	result.Addresses, err = cli.GetNewAddresses(counter)
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
