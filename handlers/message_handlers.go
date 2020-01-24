package handlers

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lodthe/cpparserbot/api"
	"github.com/lodthe/cpparserbot/buttons"
	"github.com/lodthe/cpparserbot/configs"
	"github.com/lodthe/cpparserbot/controllers"
	"github.com/lodthe/cpparserbot/helpers"
	"github.com/lodthe/cpparserbot/keyboards"
	"github.com/lodthe/cpparserbot/labels"
	"github.com/lodthe/cpparserbot/loggers"
	"github.com/lodthe/cpparserbot/models"
	"github.com/wcharczuk/go-chart"
	"strings"
	"time"
)

type MessageDispatcher struct {
	controller *controllers.TelegramController
	logger     *loggers.TelegramLogger
	binance    *api.Binance
}

//NewMessageDispatcher creates new MessageDispatcher
func NewMessageDispatcher(
	controller *controllers.TelegramController,
	logger *loggers.TelegramLogger,
	binance *api.Binance,
) *MessageDispatcher {
	return &MessageDispatcher{controller, logger, binance}
}

//DispatchMessage identifies type of message and sends response based on this type
func (d *MessageDispatcher) Dispatch(update *tgbotapi.Update) {
	text := strings.TrimSpace(update.Message.Text)

	switch true {
	//Start
	case strings.HasPrefix(text, "/start"):
		d.controller.Send(d.handleStart(update))

	//Menu
	case strings.HasPrefix(text, buttons.Menu.Text):
		d.controller.Send(d.handleMenu(update))

	//GetBinancePricesList
	case strings.HasPrefix(text, buttons.GetBinancePricesList.Text):
		d.controller.Send(d.handleGetBinancePairsList(update))

	//GetAllPrices
	case strings.HasPrefix(text, buttons.GetAllPrices.Text) || strings.HasPrefix(text, labels.GetAllCommand):
		d.controller.Send(d.handleGetAllPrices(update))

	//GetList
	case strings.HasPrefix(text, labels.GetListCommand):
		d.controller.Send(d.handleGetListCommand(update))

	//Get pair price
	case strings.HasPrefix(text, labels.GetCommand+" "):
		d.controller.Send(d.handleCommandGetBinancePrice(update))

	//Get (empty)
	case strings.HasPrefix(text, labels.GetCommand):
		d.controller.Send(d.handleGetCorrection(update))

	//GetBinancePrice
	case helpers.FindPairInConfig(text) != nil:
		d.controller.Send(d.handleGetBinancePrice(helpers.FindPairInConfig(text), update))

	default:
		d.controller.Send(d.handleUnknownCommand(update))
	}
}

//handleStart returns start message
func (d *MessageDispatcher) handleStart(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s sent /start", helpers.GetTelegramProfileURL(update)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.Start)
	msg.ReplyMarkup = keyboards.Start()
	helpers.PrepareMessageConfig(&msg)
	return msg
}

//handleUnknownCommand returns unknown command message
func (d *MessageDispatcher) handleUnknownCommand(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s sent unknown command: %s", helpers.GetTelegramProfileURL(update), update.Message.Text))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.UnknownCommand)
	msg.ReplyMarkup = keyboards.UnknownCommand()
	helpers.PrepareMessageConfig(&msg)
	return msg
}

//handleMenu returns menu message
func (d *MessageDispatcher) handleMenu(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s opened menu", helpers.GetTelegramProfileURL(update)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.Menu)
	msg.ReplyMarkup = keyboards.Menu()
	helpers.PrepareMessageConfig(&msg)
	return msg
}

//handleGetBinancePricesList returns message made of
//supported Binance pairs WITH keyboard
func (d *MessageDispatcher) handleGetBinancePairsList(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s asked Binance pairs list keyboard", helpers.GetTelegramProfileURL(update)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePairsList)
	msg.ReplyMarkup = keyboards.GetBinancePairsList()
	helpers.PrepareMessageConfig(&msg)
	return msg
}

//handleGetListCommand returns message made of
//supported Binance pairs WITH keyboard
func (d *MessageDispatcher) handleGetListCommand(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s asked Binance pairs list", helpers.GetTelegramProfileURL(update)))

	text := labels.GetList

	for _, i := range configs.BinancePairs {
		text += i.String() + "\n"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	helpers.PrepareMessageConfig(&msg)
	return msg
}

//handleGetCorrection asks user to add information
//about pair to the /get command
func (d *MessageDispatcher) handleGetCorrection(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s sent /get command without pair", helpers.GetTelegramProfileURL(update)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetCorrection)
	helpers.PrepareMessageConfig(&msg)
	return msg
}

//handleGetAllPrices return message with .xls document
//that contains all tickers information
func (d *MessageDispatcher) handleGetAllPrices(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s asked for all prices", helpers.GetTelegramProfileURL(update)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetAllPrices)
	msg.ReplyMarkup = keyboards.GetAllPrices()
	helpers.PrepareMessageConfig(&msg)
	return msg
}

//handleGetBinancePrice returns message with Binance pair price
//and a graph (which shows how price was changing during the day)
func (d *MessageDispatcher) handleGetBinancePrice(pair *models.Pair, update *tgbotapi.Update) tgbotapi.Chattable {
	profileRepr := helpers.GetTelegramProfileURL(update)
	d.logger.Info(fmt.Sprintf("%s asked for %s Binance price", profileRepr, pair))

	price, err := d.binance.GetPrice(pair)
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot get Binance price for %s (%s): %s", pair, profileRepr, err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePriceFailed)
	}

	klines, err := d.binance.GetKlines(pair)
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot get Binance klines for %s (%s): %s", pair, profileRepr, err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePriceFailed)
	}

	var x, y []float64
	for _, i := range klines {
		x = append(x, float64(i.Timestamp))
		y = append(y, i.Price)
	}

	graph := chart.Chart{
		Title: fmt.Sprintf("Изменение цены %s за сутки", pair),
		XAxis: chart.XAxis{
			ValueFormatter: func(v interface{}) string {
				t := time.Unix(int64(v.(float64))/1000, 0)
				return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
			},
		},
		YAxis: chart.YAxis{
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.8f", v.(float64))
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.ColorRed,
					FillColor:   chart.GetDefaultColor(0).WithAlpha(32),
				},
				XValues: x,
				YValues: y,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot render graph for %s (%s): %v", pair, profileRepr, err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, labels.GetBinancePriceFailed)
	}

	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "Prices", Bytes: buffer.Bytes()})
	msg.Caption = fmt.Sprintf(labels.GetBinancePrice, pair, price)
	msg.ReplyMarkup = keyboards.GetBinancePrice()
	helpers.PreparePhotoConfig(&msg)

	d.logger.Info(fmt.Sprintf("Sending message with %f price for %s to %s", price, pair, profileRepr))

	return msg
}

//handleCommandGetBinancePrice returns message with Binance pair price
//and a graph (which shows how price was changing during the day)
func (d *MessageDispatcher) handleCommandGetBinancePrice(update *tgbotapi.Update) tgbotapi.Chattable {
	//Trimming text and removing `get` command prefix
	pairName := strings.TrimSpace(update.Message.Text)[len(labels.GetCommand)+1:]
	pair := helpers.FindPairInConfig(pairName)
	if pair == nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, labels.UnknownPair)
	}

	return d.handleGetBinancePrice(pair, update)
}
