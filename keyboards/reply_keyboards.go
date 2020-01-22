package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lodthe/cpparserbot/buttons"
	"github.com/lodthe/cpparserbot/configs"
	"github.com/lodthe/cpparserbot/models"
)

//Start returns reply keyboard for `Start` message
func Start() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			buttons.Menu,
		),
	)
}

//Menu returns reply keyboard for `Menu` message
func Menu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			buttons.GetBinancePricesList,
		),
		tgbotapi.NewKeyboardButtonRow(
			buttons.GetAllPrices,
		),
	)
}

//parseButtonRow returns slice of `KeyboardButtons`, where
//buttons labels equal to pairs representation
func parseButtonRow(pairs []models.Pair) []tgbotapi.KeyboardButton {
	result := make([]tgbotapi.KeyboardButton, len(pairs))

	for _, pair := range pairs {
		result = append(result, tgbotapi.NewKeyboardButton(pair.String()))
	}

	return result
}

//min returns minimum from `a` and `b`
func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

//GetBinancePricesList returns reply keyboard for `GetBinancePricesList` message
func GetBinancePricesList() tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		buttons.Menu,
	))

	for i := 0; i < len(configs.BinancePairs); i += 2 {
		row := parseButtonRow(configs.BinancePairs[i:min(i+2, len(configs.BinancePairs))])
		rows = append(rows, row)
	}

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		buttons.GetAllPrices,
	))

	return tgbotapi.NewReplyKeyboard(rows...)
}

//GetAllPrices returns reply keyboard for `GetBinancePrices` message
func GetAllPrices() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			buttons.Menu,
		),
		tgbotapi.NewKeyboardButtonRow(
			buttons.GetBinancePricesList,
		),
	)
}

//GetBinancePrice returns reply keyboard for `GetBinancePrice` message
func GetBinancePrice() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			buttons.Menu,
		),
		tgbotapi.NewKeyboardButtonRow(
			buttons.GetAllPrices,
		),
	)
}
