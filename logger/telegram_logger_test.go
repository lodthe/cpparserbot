package logger

import (
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lodthe/cpparserbot/label"
)

const channelID int64 = 100100100

type MockMessage struct {
	ChatID int64
	Text   string
}

type MockTelegramController struct {
	Messages []MockMessage
}

func (c *MockTelegramController) Init() {}

func (c *MockTelegramController) Send(msg tgbotapi.Chattable) {
	config := msg.(*tgbotapi.MessageConfig)
	c.Messages = append(c.Messages, MockMessage{config.ChatID, config.Text})
}

func (c *MockTelegramController) Run() {}

func TestInfo(t *testing.T) {
	in := []string{"First", "Second", "Third message"}
	want := []MockMessage{
		{channelID, fmt.Sprintf(label.InfoFormat, "First")},
		{channelID, fmt.Sprintf(label.InfoFormat, "Second")},
		{channelID, fmt.Sprintf(label.InfoFormat, "Third message")},
	}

	controller := MockTelegramController{}
	logger := TelegramLogger{channelID, &controller}

	for _, i := range in {
		logger.Info(i)
	}

	for index, res := range want {
		if res != controller.Messages[index] {
			t.Fatalf("For %v expected %v, got %v", in, want, controller.Messages)
		}
	}
}

func TestError(t *testing.T) {
	in := []string{"First", "Second", "Third message"}
	want := []MockMessage{
		{channelID, fmt.Sprintf(label.ErrorFormat, "First")},
		{channelID, fmt.Sprintf(label.ErrorFormat, "Second")},
		{channelID, fmt.Sprintf(label.ErrorFormat, "Third message")},
	}

	controller := MockTelegramController{}
	logger := TelegramLogger{channelID, &controller}

	for _, i := range in {
		logger.Error(i)
	}

	for index, res := range want {
		if res != controller.Messages[index] {
			t.Fatalf("For %v expected %v, got %v", in, want, controller.Messages)
		}
	}
}
