package handlers

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lodthe/cpparserbot/api"
	"github.com/lodthe/cpparserbot/button"
	"github.com/lodthe/cpparserbot/config"
	"github.com/lodthe/cpparserbot/helper"
	"github.com/lodthe/cpparserbot/label"
	"github.com/lodthe/cpparserbot/logger"
)

func NewUpdate(ID int64, text string) *tgbotapi.Update {
	return &tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: ID,
			},
			Text: text,
		},
	}
}

type Message struct {
	ID       int64
	Text     string
	Keyboard interface{}
	FileSize int
}

func GetMessage(c tgbotapi.Chattable) *Message {
	switch c.(type) {
	case tgbotapi.MessageConfig:
		return &Message{
			ID:       c.(tgbotapi.MessageConfig).ChatID,
			Text:     c.(tgbotapi.MessageConfig).Text,
			Keyboard: c.(tgbotapi.MessageConfig).ReplyMarkup,
		}
	case *tgbotapi.MessageConfig:
		return &Message{
			ID:       c.(*tgbotapi.MessageConfig).ChatID,
			Text:     c.(*tgbotapi.MessageConfig).Text,
			Keyboard: c.(*tgbotapi.MessageConfig).ReplyMarkup,
		}
	case *tgbotapi.DocumentConfig:
		return &Message{
			ID:       c.(*tgbotapi.DocumentConfig).ChatID,
			Text:     c.(*tgbotapi.DocumentConfig).Caption,
			Keyboard: c.(*tgbotapi.DocumentConfig).ReplyMarkup,
			FileSize: c.(*tgbotapi.DocumentConfig).FileSize,
		}
	case *tgbotapi.PhotoConfig:
		return &Message{
			ID:       c.(*tgbotapi.PhotoConfig).ChatID,
			Text:     c.(*tgbotapi.PhotoConfig).Caption,
			Keyboard: c.(*tgbotapi.PhotoConfig).ReplyMarkup,
			FileSize: c.(*tgbotapi.PhotoConfig).FileSize,
		}
	default:
		log.Println(fmt.Sprintf("UNKNOWN CONFIG: %T", c))
		return nil
	}
}

type MockTelegramController struct {
	mu       sync.Mutex
	Messages []*Message
}

func (c *MockTelegramController) Init() {
	c.mu = sync.Mutex{}
}

func (c *MockTelegramController) Send(msg tgbotapi.Chattable) {
	c.mu.Lock()
	c.Messages = append(c.Messages, GetMessage(msg))
	c.mu.Unlock()
}

func (c *MockTelegramController) Run() {}

func NewControllerAndDispatcher() (*MockTelegramController, *MessageDispatcher) {
	c := &MockTelegramController{}
	c.Init()
	fakeC := &MockTelegramController{}
	fakeC.Init()
	b := &api.Binance{}
	b.Init(os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_SECRET_KEY"))
	return c, NewMessageDispatcher(c, &logger.TelegramLogger{Controller: fakeC}, b)
}

func TestStartMessage(t *testing.T) {
	in := []*tgbotapi.Update{
		NewUpdate(-1, "/start"),
		NewUpdate(1, "/start"),
		NewUpdate(-1, "UNKNOWN COMMAND"),
	}

	controller, dispatcher := NewControllerAndDispatcher()

	for _, u := range in {
		dispatcher.Dispatch(u)
	}

	for i := range in {
		keyboard := controller.Messages[i].Keyboard
		if (keyboard != nil) != (helper.GetChatID(in[i]) > 0) {
			t.Errorf("Expected chat existance status %v, got %v: %v", helper.GetChatID(in[i]) >= 0, keyboard != nil, keyboard)
		}

		if (in[i].Message.Text == "/start") && (controller.Messages[i].Text != label.Start) {
			t.Errorf("Expected message text \n%s, \ngot \n%s", label.Start, controller.Messages[i].Text)
		}
	}
}

// Test validates keyboard removing and equality of
// getting pair price commands (with and without /get prefix)
func TestGetPrice(t *testing.T) {
	pairs := config.BinancePairs

	// GetPrice with /get prefix
	c1, d1 := NewControllerAndDispatcher()
	// GetPrice without /get prefix
	c2, d2 := NewControllerAndDispatcher()
	for i, pair := range pairs {
		ID := int64(1)
		if i%2 == 0 {
			ID = -1
		}

		d1.Dispatch(NewUpdate(ID, fmt.Sprintf("%s %s", label.GetCommand, pair)))
		d2.Dispatch(NewUpdate(ID, pair.String()))
	}

	for i := range pairs {
		if len(c1.Messages[i].Text) != len(c2.Messages[i].Text) {
			t.Errorf("Messages length are not equal: \n%s\n\n%s", c1.Messages[i].Text, c2.Messages[i].Text)
		}

		if c1.Messages[i].FileSize != c2.Messages[i].FileSize {
			t.Errorf("FileSizes are not equal: \n%d\n\n%d", c1.Messages[i].FileSize, c2.Messages[i].FileSize)
		}
	}
}

func TestGetAllBinancePrices(t *testing.T) {
	in := []*tgbotapi.Update{
		NewUpdate(-1, label.GetAllBinanceCommand),
		NewUpdate(1, label.GetAllBinanceCommand),
	}

	controller, dispatcher := NewControllerAndDispatcher()

	for _, u := range in {
		dispatcher.Dispatch(u)
	}

	for i := range in {
		keyboard := controller.Messages[i].Keyboard
		if (keyboard != nil) != (helper.GetChatID(in[i]) > 0) {
			t.Errorf("Expected chat existance status %v, got %v: %v", helper.GetChatID(in[i]) >= 0, keyboard != nil, keyboard)
		}

		if controller.Messages[i].FileSize == 0 {
			t.Errorf("Expected document, got file size == 0: %v", controller.Messages[i])
		}
	}
}

func TestGetAllPrices(t *testing.T) {
	in := []*tgbotapi.Update{
		NewUpdate(-1, label.GetAllCommand),
		NewUpdate(1, label.GetAllCommand),
		NewUpdate(-1, button.GetAllPrices.Text),
		NewUpdate(1, button.GetAllPrices.Text),
	}

	controller, dispatcher := NewControllerAndDispatcher()

	for _, u := range in {
		dispatcher.Dispatch(u)
	}

	for i := range in {
		keyboard := controller.Messages[i].Keyboard
		if (keyboard != nil) != (helper.GetChatID(in[i]) > 0) {
			t.Errorf("Expected chat existance status %v, got %v: %v", helper.GetChatID(in[i]) >= 0, keyboard != nil, keyboard)
		}

		if controller.Messages[i].FileSize == 0 {
			t.Errorf("Expected document, got file size == 0: %v", controller.Messages[i])
		}
	}
}
