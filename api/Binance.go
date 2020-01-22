package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/lodthe/ratesparser/models"
)

const (
	binanceAPIURL = "https://api.binance.com/api/v3/"
)

//Binance provides with binance.com API methods
type Binance struct{ Counter int }

//BinanceTicker keeps information about binance ticker
type binanceTicker struct {
	Symbol string
	Price  string
}

func (exchanger *Binance) GetPrice(pair models.Pair) (float64, error) {
	query := binanceAPIURL + "ticker/price?symbol=" + pair.SpendCurrency + pair.BuyCurrency
	response, err := http.Get(query)
	if err != nil {
		return 0, fmt.Errorf("cannot get %s price from Binance", pair)
	}
	defer response.Body.Close()

	var ticker binanceTicker
	body, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(body, &ticker)

	result, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse price for %v from Binance response", pair)
	}

	return result, nil
}
