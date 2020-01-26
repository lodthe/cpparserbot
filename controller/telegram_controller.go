// Package controller implements controllers such as TelegramController that manages message sending
package controller

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const maxMessagesCountPerSecond = 30

// Controller interface describes TelegramController methods for tests
type Controller interface {
	Init()
	Send(msg tgbotapi.Chattable)
	Run()
}

// BotAPI interface is made to have an opportunity of mocking tgbotapi.BotAPI
type BotAPI interface {
	Send(msg tgbotapi.Chattable) (tgbotapi.Message, error)
}

// TelegramController controls request sending to
// prevent against exceeding Telegram limits
type TelegramController struct {
	Bot          BotAPI
	messages     chan tgbotapi.Chattable
	sentMessages []time.Time
}

// Init initializes controller:
//  sets messages channel buffer size
func (controller *TelegramController) Init() {
	controller.messages = make(chan tgbotapi.Chattable, maxMessagesCountPerSecond*2)
}

// Send pushes message to the channel with messages,
// which will be send in the future (Run controls this)
func (controller *TelegramController) Send(msg tgbotapi.Chattable) {
	go func() {
		controller.messages <- msg
	}()
}

// Run starts sending messages to telegram and controls limits:
// no more then `maxMessagesCountPerSecond`
func (controller *TelegramController) Run() {
	go func(controller *TelegramController) {
		for {
			// Removing messages, sent more than 1 second ago, from the queue
			for (len(controller.sentMessages) != 0) && (time.Now().Sub(controller.sentMessages[0]) > time.Second) {
				controller.sentMessages = controller.sentMessages[1:]
			}

			if len(controller.sentMessages) < maxMessagesCountPerSecond {
				// Sending message
				go func(c tgbotapi.Chattable) {
					controller.Bot.Send(c)
				}(<-controller.messages)

				controller.sentMessages = append(controller.sentMessages, time.Now())
			} else {
				time.Sleep(time.Second / maxMessagesCountPerSecond)
			}
		}
	}(controller)
}
