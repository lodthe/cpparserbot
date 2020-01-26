// Package api implements Binance api methods
package api

import (
	"context"
	"strconv"
	"time"

	"github.com/adshao/go-binance"

	"github.com/lodthe/cpparserbot/model"
)

// Binance provides with binance.com API methods
type Binance struct {
	client *binance.Client
}

// Init initializes Binance client with API keys
func (b *Binance) Init(apiKey, secretKey string) {
	b.client = binance.NewClient(apiKey, secretKey)
}

// GetPrice returns Binance price for given pair
func (b *Binance) GetPrice(pair *model.Pair) (float64, error) {
	prices, err := b.client.NewListPricesService().Symbol(pair.ToBinanceFormat()).Do(context.Background())
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(prices[0].Price, 64)
}

// GetAllPrices returns data with Binance prices for all pairs
func (b *Binance) GetAllPrices() ([]*binance.SymbolPrice, error) {
	return b.client.NewListPricesService().Do(context.Background())
}

// Kline holds information about Binance kline record
type Kline struct {
	Price     float64
	Timestamp int64
}

// GetKlines returns information about how pair price was changing during the day
func (b *Binance) GetKlines(pair *model.Pair) ([]Kline, error) {
	klines, err := b.client.
		NewKlinesService().Symbol(pair.ToBinanceFormat()).
		Interval("1h").
		StartTime(int64(1000) * (time.Now().Add(-time.Hour * 24).Unix())).
		Do(context.Background())
	if err != nil {
		return make([]Kline, 0), err
	}

	var result []Kline

	// Extracting data from response
	for _, i := range klines {
		price, err := strconv.ParseFloat(i.Close, 64)
		if err != nil {
			return make([]Kline, 0), err
		}

		result = append(result, Kline{price, i.CloseTime})
	}

	return result, nil
}
