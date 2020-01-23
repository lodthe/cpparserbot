package models

import (
	"fmt"
)

//Pair keeps information about currency pair
type Pair struct {
	SpendCurrency string
	BuyCurrency   string
}

//String returns Pair representation as string
func (pair Pair) String() string {
	return fmt.Sprintf("%s/%s", pair.SpendCurrency, pair.BuyCurrency)
}

//ToBinanceFormat returns Pair as Binance symbol
func (pair *Pair) ToBinanceFormat() string {
	return fmt.Sprintf("%s%s", pair.SpendCurrency, pair.BuyCurrency)
}
