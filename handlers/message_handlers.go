package handlers

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lodthe/cpparserbot/api"
	"github.com/lodthe/cpparserbot/buttons"
	"github.com/lodthe/cpparserbot/controllers"
	"github.com/lodthe/cpparserbot/helpers"
	"github.com/lodthe/cpparserbot/keyboards"
	"github.com/lodthe/cpparserbot/labels"
	"github.com/lodthe/cpparserbot/loggers"
	"github.com/wcharczuk/go-chart"
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
func handleStart(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.Start)
	msg.ReplyMarkup = keyboards.Start()
	return msg
}

//handleStart returns unknown command message
func handleUnknownCommand(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.UnknownCommand)
	return msg
}

//handleMenu returns menu message
func handleMenu(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.Menu)
	msg.ReplyMarkup = keyboards.Menu()
	return msg
}

//handleGetBinancePricesList returns message with keyboard
//that consists of available Binance tickers
func handleGetBinancePricesList(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePricesList)
	msg.ReplyMarkup = keyboards.GetBinancePricesList()
	return msg
}

//handleGetAllPrices return message with .xls document
//that contains all tickers information
func handleGetAllPrices(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetAllPrices)
	return msg
}

//handleGetBinancePrice return message with Binance pair price
//and a graph (which shows how price was changing during the day)
func handleGetBinancePrice(update tgbotapi.Update, binanceAPI *api.Binance) tgbotapi.Chattable {
	pair := *helpers.FindPairInConfig(update.Message.Text)
	price, err := binanceAPI.GetPrice(pair)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePriceFailed)
	}

	klines, err := binanceAPI.GetKlines(pair)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePriceFailed)
	}

	var x, y []float64
	for _, i := range klines {
		x = append(x, float64(i.Timestamp))
		y = append(y, i.Price)
	}

	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: x,
				YValues: y,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePriceFailed)
	}

	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "Prices", Bytes: buffer.Bytes()})
	msg.Caption = fmt.Sprintf(labels.GetBinancePrice, pair, price)

	return msg
}
