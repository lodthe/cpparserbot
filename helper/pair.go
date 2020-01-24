// Package helper implements different helpful functions
package helper

import (
	"github.com/lodthe/cpparserbot/config"
	"github.com/lodthe/cpparserbot/model"
	"strings"
)

// findPairInConfig checks if `s` occurs in Binance pairs defined in config
func FindPairInConfig(s string) *model.Pair {
	for _, pair := range config.BinancePairs {
		if strings.EqualFold(s, pair.String()) {
			return &pair
		}
	}
	return nil
}
