// Package config keeps config information such as Binance pairs list
package config

import "github.com/lodthe/cpparserbot/model"

var (
	// Supported Binance pairs
	BinancePairs = []model.Pair{
		{"BTC", "USDT"},
		{"BNB", "BTC"},
		{"ETH", "BTC"},
		{"XRP", "BTC"},
		{"DASH", "BTC"},
		{"ETC", "BTC"},
	}
)
