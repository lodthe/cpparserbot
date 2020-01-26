package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Rate keeps rate information
type Rate struct {
	Pair      string
	Exchange  string
	BuyPrice  string
	SellPrice string
	Updated   string
}

// GetRates returns information about all rates
func GetRates() ([]Rate, error) {
	response, err := http.Get("https://frates.herokuapp.com/get_rates")
	if err != nil {
		return make([]Rate, 0), errors.New("cannot get all rates")
	}
	defer response.Body.Close()

	var rates []Rate
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return make([]Rate, 0), errors.New("cannot unmarshal response")
	}
	return rates, nil
}
