package daemon

import "encoding/json"

type ListAddressGroupingsAddress struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
	Label   string  `json:"label"`
}

type ListAddressGroupingsResponse struct {
	Addresses []ListAddressGroupingsAddress `json:"addresses"`
}

func (r *ListAddressGroupingsResponse) UnmarshalJSON(data []byte) error {
	var rawAddressGroups [][][]interface{}
	err := json.Unmarshal(data, &rawAddressGroups)
	if err != nil {
		return err
	}

	for _, group := range rawAddressGroups {
		for _, item := range group {
			address := item[0].(string)
			amount := item[1].(float64)
			var label string
			if len(item) > 2 {
				label = item[2].(string)
			}

			r.Addresses = append(r.Addresses, ListAddressGroupingsAddress{
				Address: address,
				Amount:  amount,
				Label:   label,
			})
		}
	}

	return nil
}
