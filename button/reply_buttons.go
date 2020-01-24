// Package button keeps keyboard buttons
package button

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var (
	Menu                 = tgbotapi.NewKeyboardButton("Меню")
	GetBinancePricesList = tgbotapi.NewKeyboardButton("Узнать курс на Binance")
	GetAllPrices         = tgbotapi.NewKeyboardButton("Получить все курсы")
)
