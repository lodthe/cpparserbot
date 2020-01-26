// Package helper implements different helpful functions
package helper

import (
	"strings"

	"github.com/lodthe/cpparserbot/config"
	"github.com/lodthe/cpparserbot/model"
)

var binancePairs []model.Pair

func init() {
	binancePairs = config.BinancePairs
}

// findPairInConfig checks whether `s` occurs in Binance pairs defined in config
func FindPairInConfig(s string) *model.Pair {
	for _, pair := range binancePairs {
		if strings.EqualFold(s, pair.String()) {
			return &pair
		}
	}
	return nil
}
