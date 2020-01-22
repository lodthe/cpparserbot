package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lodthe/cpparserbot/buttons"
	"github.com/lodthe/cpparserbot/controllers"
	"github.com/lodthe/cpparserbot/keyboards"
	"github.com/lodthe/cpparserbot/labels"
	"github.com/lodthe/cpparserbot/loggers"
)

//DispatchMessage identifies type of message and sends response based on this type
func DispatchMessage(update tgbotapi.Update, controller *controllers.TelegramController, logger *loggers.TelegramLogger) {
	chatID := update.Message.Chat.ID

	switch update.Message.Text {
	case "/start":
		logger.Info(fmt.Sprintf("[%v](tg://user?id=%v) sent /start", chatID, chatID))
		controller.Send(handleStart(update))

	case buttons.Menu.Text:
		controller.Send(handleMenu(update))

	case buttons.GetBinancePricesList.Text:
		controller.Send(handleGetBinancePricesList(update))

	case buttons.GetAllPrices.Text:
		controller.Send(handleGetAllPrices(update))
	}
}

//handleStart returns start message
func handleStart(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.Start)
	msg.ReplyMarkup = keyboards.Start()
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
	msg.ReplyMarkup = keyboards.GetAllPrices()
	return msg
}
