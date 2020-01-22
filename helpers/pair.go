package helpers

import (
	"github.com/lodthe/cpparserbot/configs"
	"github.com/lodthe/cpparserbot/models"
)

//findPairInConfig checks if `s` occurs in Binance pairs defined in config
func FindPairInConfig(s string) *models.Pair {
	for _, pair := range configs.BinancePairs {
		if s == pair.String() {
			return &pair
		}
	}
	return nil
}
