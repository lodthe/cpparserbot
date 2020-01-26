// Package handlers implements user request handling
package handlers

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tealeg/xlsx"

	"github.com/lodthe/cpparserbot/api"
	"github.com/lodthe/cpparserbot/button"
	"github.com/lodthe/cpparserbot/config"
	"github.com/lodthe/cpparserbot/controller"
	"github.com/lodthe/cpparserbot/helper"
	"github.com/lodthe/cpparserbot/keyboard"
	"github.com/lodthe/cpparserbot/label"
	"github.com/lodthe/cpparserbot/logger"
	"github.com/lodthe/cpparserbot/model"
	"github.com/wcharczuk/go-chart"
)

// MessageDispatcher holds additional fields for handling messages
type MessageDispatcher struct {
	controller controller.Controller
	logger     *logger.TelegramLogger
	binance    *api.Binance
}

// NewMessageDispatcher creates new MessageDispatcher
func NewMessageDispatcher(
	controller controller.Controller,
	logger *logger.TelegramLogger,
	binance *api.Binance,
) *MessageDispatcher {
	return &MessageDispatcher{controller, logger, binance}
}

// DispatchMessage identifies type of message and sends response based on this type
func (d *MessageDispatcher) Dispatch(update *tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			d.logger.Error(fmt.Sprintf("Panic! Cannot dispatch message\n%v\nfrom %v\n\n%v",
				update.Message.Text, update.Message.Chat.ID, r))
		}
	}()
	text := strings.TrimSpace(update.Message.Text)

	switch true {
	// Start
	case strings.HasPrefix(text, "/start"):
		d.controller.Send(d.handleStart(update))

	// Menu
	case strings.HasPrefix(text, button.Menu.Text):
		d.controller.Send(d.handleMenu(update))

	// GetBinancePricesList
	case strings.HasPrefix(text, button.GetBinancePricesList.Text):
		d.controller.Send(d.handleGetBinancePairsList(update))

	// GetAllBinancePrices
	case strings.HasPrefix(text, label.GetAllBinanceCommand):
		d.controller.Send(d.handleGetAllBinancePrices(update))

	// GetAllPrices
	case strings.HasPrefix(text, label.GetAllCommand) || strings.HasPrefix(text, button.GetAllPrices.Text):
		d.controller.Send(d.handleGetAllPrices(update))

	// GetList
	case strings.HasPrefix(text, label.GetListCommand):
		d.controller.Send(d.handleGetListCommand(update))

	// Get pair price
	case strings.HasPrefix(text, label.GetCommand+" "):
		d.controller.Send(d.handleCommandGetBinancePrice(update))

	// Get (empty)
	case strings.HasPrefix(text, label.GetCommand):
		d.controller.Send(d.handleGetCorrection(update))

	// GetBinancePrice
	case helper.FindPairInConfig(text) != nil:
		d.controller.Send(d.handleGetBinancePrice(helper.FindPairInConfig(text), update))

	default:
		d.controller.Send(d.handleUnknownCommand(update))
	}
}

// handleStart returns start message
func (d *MessageDispatcher) handleStart(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s sent /start", helper.GetTelegramProfileURL(update)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, label.Start)
	msg.ReplyMarkup = keyboard.Start()
	helper.PrepareMessage(&msg)
	return &msg
}

// handleUnknownCommand returns unknown command message
func (d *MessageDispatcher) handleUnknownCommand(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s sent unknown command: %s", helper.GetTelegramProfileURL(update), update.Message.Text))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, label.UnknownCommand)
	msg.ReplyMarkup = keyboard.UnknownCommand()
	helper.PrepareMessage(&msg)
	return &msg
}

// handleMenu returns menu message
func (d *MessageDispatcher) handleMenu(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s opened menu", helper.GetTelegramProfileURL(update)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, label.Menu)
	msg.ReplyMarkup = keyboard.Menu()
	helper.PrepareMessage(&msg)
	return &msg
}

// handleGetBinancePricesList returns message made of
// supported Binance pairs WITH keyboard
func (d *MessageDispatcher) handleGetBinancePairsList(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s asked Binance pairs list keyboard", helper.GetTelegramProfileURL(update)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, label.GetBinancePairsList)
	msg.ReplyMarkup = keyboard.GetBinancePairsList()
	helper.PrepareMessage(&msg)
	return &msg
}

// handleGetListCommand returns message made of
// supported Binance pairs WITH keyboard
func (d *MessageDispatcher) handleGetListCommand(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s asked Binance pairs list", helper.GetTelegramProfileURL(update)))

	text := label.GetList

	for _, i := range config.BinancePairs {
		text += i.String() + "\n"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	helper.PrepareMessage(&msg)
	return &msg
}

// handleGetCorrection asks user to add information
// about pair to the /get command
func (d *MessageDispatcher) handleGetCorrection(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s sent /get command without pair", helper.GetTelegramProfileURL(update)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, label.GetCorrection)
	helper.PrepareMessage(&msg)
	return &msg
}

// handleGetAllBinancePrices return message with .xls document
// that contains all Binance tickers information
func (d *MessageDispatcher) handleGetAllBinancePrices(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s asked for all binance prices", helper.GetTelegramProfileURL(update)))

	prices, err := d.binance.GetAllPrices()
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot get all Binance prices: %s", err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.GetAllBinancePricesFailed)
	}

	// Creating xlsx table
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Binance")
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot add sheet to the Binance prices table: %s", err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.GetAllBinancePricesFailed)
	}

	header := sheet.AddRow()
	header.AddCell().Value = "Symbol"
	header.AddCell().Value = "Price"

	for _, price := range prices {
		sheet.AddRow().WriteStruct(price, 2)
	}

	var buffer bytes.Buffer
	err = file.Write(&buffer)
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot save XLSX table with Binance prices: %s", err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.GetAllBinancePricesFailed)
	}

	// Preparing message
	msg := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "Binance.xlsx", Bytes: buffer.Bytes()})
	msg.FileSize = len(buffer.Bytes())
	msg.Caption = label.GetAllBinancePrices
	msg.ReplyMarkup = keyboard.GetAllBinancePrices()
	helper.PrepareMessage(&msg)

	d.logger.Info(fmt.Sprintf("Sending message with all Binance prices to %s", helper.GetTelegramProfileURL(update)))

	return &msg
}

// handleGetAllPrices return message with .xls document
// that contains all tickers information
func (d *MessageDispatcher) handleGetAllPrices(update *tgbotapi.Update) tgbotapi.Chattable {
	d.logger.Info(fmt.Sprintf("%s asked for all prices", helper.GetTelegramProfileURL(update)))

	prices, err := api.GetRates()
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot get all prices: %s", err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.GetAllPricesFailed)
	}

	// Creating xlsx table
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Prices")
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot add sheet to the prices table: %s", err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.GetAllPricesFailed)
	}

	header := sheet.AddRow()
	val := reflect.ValueOf(&api.Rate{}).Elem()

	for i := 0; i < val.NumField(); i++ {
		header.AddCell().Value = val.Type().Field(i).Name
	}

	for _, price := range prices {
		sheet.AddRow().WriteStruct(&price, val.NumField())
	}

	var buffer bytes.Buffer
	err = file.Write(&buffer)
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot save XLSX table with all prices: %s", err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.GetAllPricesFailed)
	}

	// Preparing message
	msg := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "Prices.xlsx", Bytes: buffer.Bytes()})
	msg.FileSize = len(buffer.Bytes())
	msg.Caption = label.GetAllPrices
	msg.ReplyMarkup = keyboard.GetAllPrices()
	helper.PrepareMessage(&msg)

	d.logger.Info(fmt.Sprintf("Sending message with all prices to %s", helper.GetTelegramProfileURL(update)))

	return &msg
}

// handleGetBinancePrice returns message with Binance pair price
// and a graph (which shows how price was changing during the day)
func (d *MessageDispatcher) handleGetBinancePrice(pair *model.Pair, update *tgbotapi.Update) tgbotapi.Chattable {
	profileRepr := helper.GetTelegramProfileURL(update)
	d.logger.Info(fmt.Sprintf("%s asked for %s Binance price", profileRepr, pair))

	// Getting pair price
	price, err := d.binance.GetPrice(pair)
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot get Binance price for %s (%s): %s", pair, profileRepr, err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.GetBinancePriceFailed)
	}

	// Getting klines
	klines, err := d.binance.GetKlines(pair)
	if err != nil {
		d.logger.Error(fmt.Sprintf("Cannot get Binance klines for %s (%s): %s", pair, profileRepr, err))
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.GetBinancePriceFailed)
	}

	var x, y []float64
	for _, i := range klines {
		x = append(x, float64(i.Timestamp))
		y = append(y, i.Price)
	}

	// Creating chart with klines
	graph := chart.Chart{
		Title: fmt.Sprintf("Изменение цены %s за сутки", pair),
		XAxis: chart.XAxis{
			// Convert UNIX-time to user-readable format
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
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.GetBinancePriceFailed)
	}

	// Preparing message
	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "Prices.png", Bytes: buffer.Bytes()})
	msg.Caption = fmt.Sprintf(label.GetBinancePrice, pair, price)
	msg.ReplyMarkup = keyboard.GetBinancePrice()
	helper.PrepareMessage(&msg)

	d.logger.Info(fmt.Sprintf("Sending message with %f price for %s to %s", price, pair, profileRepr))

	return &msg
}

// handleCommandGetBinancePrice returns message with Binance pair price
// and a graph (which shows how price was changing throughout the day)
func (d *MessageDispatcher) handleCommandGetBinancePrice(update *tgbotapi.Update) tgbotapi.Chattable {
	// Trimming text and removing `get` command prefix
	pairName := strings.TrimSpace(update.Message.Text)[len(label.GetCommand)+1:]
	pair := helper.FindPairInConfig(pairName)
	if pair == nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, label.UnknownPair)
	}

	return d.handleGetBinancePrice(pair, update)
}
