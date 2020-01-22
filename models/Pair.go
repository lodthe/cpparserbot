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
func (pair *Pair) String() string {
	return fmt.Sprintf("%s/%s", pair.SpendCurrency, pair.BuyCurrency)
}
