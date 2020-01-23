package api

import (
	"context"
	"strconv"

	"github.com/adshao/go-binance"

	"github.com/lodthe/cpparserbot/models"
)

//Binance provides with binance.com API methods
type Binance struct {
	client *binance.Client
}

//Init initializes Binance client with API keys
func (b *Binance) Init(apiKey, secretKey string) {
	b.client = binance.NewClient(apiKey, secretKey)
}

//GetPrice returns Binance price for given pair
func (b *Binance) GetPrice(pair models.Pair) (float64, error) {
	prices, err := b.client.NewListPricesService().Symbol(pair.ToBinanceFormat()).Do(context.Background())
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(prices[0].Price, 64)
}
