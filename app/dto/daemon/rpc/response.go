package daemon

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/dto/daemon"
)

type BaseResponse struct {
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
	Id     string      `json:"id"`
}

type ListAddressGroupingsResponse struct {
	Result []daemon.ListAddressGroupingsAddress `json:"result"`
	Error  interface{}                          `json:"error"`
	Id     string                               `json:"id"`
}

func (r *ListAddressGroupingsResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	// Extract the result from the map
	result, ok := raw["result"]
	if !ok {
		return nil
	}

	// Type assertion to []interface{} to work with the nested data
	topLevelArray, ok := result.([]interface{})
	if !ok {
		log.Warn("ListAddressGroupingsResponse contains invalid Result key")
		return nil
	}

	// Process the top-level array
	for _, groupInterface := range topLevelArray {
		group, ok := groupInterface.([]interface{})
		if !ok {
			log.Warn("ListAddressGroupingsResponse contains an invalid group")
			continue
		}

		// Prepare to collect the individual address groupings
		var addressGroup []daemon.ListAddressGroupingsAddress

		// Process the second-level array
		for _, itemInterface := range group {
			item, ok := itemInterface.([]interface{})
			if !ok || len(item) < 2 { // Ensure there are at least two elements (address and amount)
				log.Warn("ListAddressGroupingsResponse contains an invalid item")
				continue
			}

			// Extract and type assert the address
			address, ok := item[0].(string)
			if !ok {
				log.Warn("Invalid address format")
				continue
			}

			// Extract and type assert the amount
			amount, ok := item[1].(float64)
			if !ok {
				log.Warn("Invalid amount format")
				continue
			}

			// Extract the label if present
			var label string
			if len(item) > 2 {
				label, _ = item[2].(string)
			}

			// Append the parsed result to the result list
			addressGroup = append(addressGroup, daemon.ListAddressGroupingsAddress{
				Address: address,
				Amount:  amount,
				Label:   label,
			})
		}

		// Append the address group to the overall result
		r.Result = append(r.Result, addressGroup...)
	}

	return nil
}

type ListUnspentResponse struct {
	Result []daemon.Unspent `json:"result"`
	Error  interface{}      `json:"error"`
	Id     string           `json:"id"`
}
