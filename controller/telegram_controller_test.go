package controller

import (
	"sort"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const testMessagesCount = 2 * maxMessagesCountPerSecond

// Highest allowed deviation of sent messages per second
// from maxMessagesCountPerSecond rate
const tolerance = 1.15

type MockBotAPI struct {
	Messages chan tgbotapi.Chattable
}

func (b *MockBotAPI) Init() {
	b.Messages = make(chan tgbotapi.Chattable, 0)
}

func (b *MockBotAPI) Send(msg tgbotapi.Chattable) (tgbotapi.Message, error) {
	go func(b *MockBotAPI, c tgbotapi.Chattable) {
		b.Messages <- c
	}(b, msg)

	return tgbotapi.Message{}, nil
}

func TestMessageDeliveryAndTimeLimit(t *testing.T) {
	in := make([]string, 0)
	for i := 0; i < testMessagesCount; i++ {
		in = append(in, string(i))
	}

	bot := &MockBotAPI{}
	bot.Init()
	controller := TelegramController{Bot: bot}
	controller.Init()
	go controller.Run()

	// Sending all messages
	for i := range in {
		controller.Send(tgbotapi.PhotoConfig{
			Caption: in[i],
		})
	}

	response := make([]string, 0)

	// Starting collecting messages
	go func() {
		for i := 0; i < testMessagesCount; i++ {
			response = append(response, (<-bot.Messages).(tgbotapi.PhotoConfig).Caption)
		}
	}()

	// Checking whether we receive not more requests per second than allowed
	for i := 0; i < testMessagesCount/maxMessagesCountPerSecond; i++ {
		count := len(response)
		time.Sleep(time.Second)
		sent := len(response)

		if float64(count) > tolerance*maxMessagesCountPerSecond {
			t.Errorf("Expected <= %d requests, got %d", maxMessagesCountPerSecond, sent)
		}
	}

	// Validating received information
	if len(response) != testMessagesCount {
		t.Errorf("Expected %d requests, got %d", testMessagesCount, len(response))
	}

	sort.Strings(response)
	sort.Strings(in)

	for i, r := range response {
		if r != in[i] {
			t.Errorf("Expected %v, got %v", in, response)
		}
	}
}
