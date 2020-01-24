// Package keyboard implements Telegram keyboards
package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lodthe/cpparserbot/button"
	"github.com/lodthe/cpparserbot/config"
	"github.com/lodthe/cpparserbot/model"
)

// Start returns reply keyboard for `Start` message
func Start() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			button.Menu,
		),
	)
}

// UnknownCommand returns reply keyboard for `UnknownCommand` message
func UnknownCommand() tgbotapi.ReplyKeyboardMarkup {
	return Menu()
}

// Menu returns reply keyboard for `Menu` message
func Menu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			button.GetBinancePricesList,
		),
		tgbotapi.NewKeyboardButtonRow(
			button.GetAllPrices,
		),
	)
}

// GetBinancePrice returns reply keyboard for `GetBinancePrice` message
func GetBinancePrice() tgbotapi.ReplyKeyboardMarkup {
	return GetBinancePairsList()
}

// parseButtonRow returns slice of `KeyboardButtons`, where
// button label equal to pairs representation
func parseButtonRow(pairs []model.Pair) []tgbotapi.KeyboardButton {
	result := make([]tgbotapi.KeyboardButton, len(pairs))

	for _, pair := range pairs {
		result = append(result, tgbotapi.NewKeyboardButton(pair.String()))
	}

	return result
}

// min returns low for `a` and `b`
func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

// GetBinancePairsList returns reply keyboard for `GetBinancePairsList` message
func GetBinancePairsList() tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		button.Menu,
	))

	for i := 0; i < len(config.BinancePairs); i += 2 {
		row := parseButtonRow(config.BinancePairs[i:min(i+2, len(config.BinancePairs))])
		rows = append(rows, row)
	}

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		button.GetAllPrices,
	))

	return tgbotapi.NewReplyKeyboard(rows...)
}

// GetAllPrices returns reply keyboard for `GetAllPrices` message
func GetAllPrices() tgbotapi.ReplyKeyboardMarkup {
	return Menu()
}
