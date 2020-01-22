package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lodthe/cpparserbot/buttons"
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

//GetBinancePricesList returns reply keyboard for `GetBinancePricesList` message
func GetBinancePricesList() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			buttons.Menu,
		),
		tgbotapi.NewKeyboardButtonRow(
			buttons.GetAllPrices,
		),
	)
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
