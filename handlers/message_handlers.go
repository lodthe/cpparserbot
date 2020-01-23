package handlers

import (
	"fmt"
	"github.com/lodthe/cpparserbot/api"
	"github.com/lodthe/cpparserbot/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lodthe/cpparserbot/buttons"
	"github.com/lodthe/cpparserbot/controllers"
	"github.com/lodthe/cpparserbot/keyboards"
	"github.com/lodthe/cpparserbot/labels"
	"github.com/lodthe/cpparserbot/loggers"
)

//DispatchMessage identifies type of message and sends response based on this type
func DispatchMessage(
	update tgbotapi.Update,
	controller *controllers.TelegramController,
	logger *loggers.TelegramLogger,
	binanceAPI *api.Binance,
) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	switch true {
	case update.Message.Text == "/start":
		logger.Info(fmt.Sprintf("[%v](tg://user?id=%v) sent /start", chatID, chatID))
		controller.Send(handleStart(update))

	case update.Message.Text == buttons.Menu.Text:
		controller.Send(handleMenu(update))

	case update.Message.Text == buttons.GetBinancePricesList.Text:
		controller.Send(handleGetBinancePricesList(update))

	case update.Message.Text == buttons.GetAllPrices.Text:
		controller.Send(handleGetAllPrices(update))

	case helpers.FindPairInConfig(text) != nil:
		logger.Info(fmt.Sprintf("[%v](tg://user?id=%v) asked for %s Binance price", chatID, chatID, text))
		controller.Send(handleGetBinancePrice(update, binanceAPI))

	default:
		controller.Send(handleUnknownCommand(update))
	}
}

//handleStart returns start message
func handleStart(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.Start)
	msg.ReplyMarkup = keyboards.Start()
	return msg
}

//handleStart returns unknown command message
func handleUnknownCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.UnknownCommand)
	return msg
}

//handleMenu returns menu message
func handleMenu(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.Menu)
	msg.ReplyMarkup = keyboards.Menu()
	return msg
}

//handleGetBinancePricesList returns message with keyboard
//that consists of available Binance tickers
func handleGetBinancePricesList(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePricesList)
	msg.ReplyMarkup = keyboards.GetBinancePricesList()
	return msg
}

//handleGetAllPrices return message with .xls document
//that contains all tickers information
func handleGetAllPrices(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetAllPrices)
	return msg
}

//handleGetAllPrices return message with .xls document
//that contains all tickers information
func handleGetBinancePrice(update tgbotapi.Update, binanceAPI *api.Binance) tgbotapi.MessageConfig {
	pair := *helpers.FindPairInConfig(update.Message.Text)
	price, err := binanceAPI.GetPrice(pair)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePriceFailed)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(labels.GetBinancePrice, pair, price))
	return msg
}
