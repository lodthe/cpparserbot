package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

const maxMessagesCountPerSecond = 25

//TelegramController controls request sending to
//prevent against exceeding Telegram limits
type TelegramController struct {
	Bot          *tgbotapi.BotAPI
	messages     chan tgbotapi.Chattable
	sentMessages []time.Time
}

//Init initializes controller:
//(*) sets messages channel buffer size
func (controller *TelegramController) Init() {
	controller.messages = make(chan tgbotapi.Chattable, maxMessagesCountPerSecond*2)
}

//Send pushes message to the channel with messages,
//which will be send in the future (Run controls this)
func (controller *TelegramController) Send(msg tgbotapi.Chattable) {
	go func() {
		controller.messages <- msg
	}()
}

//Run starts sending messages to telegram and controls limits:
//no more then `maxMessagesCountPerSecond`
func (controller *TelegramController) Run() {
	go func(controller *TelegramController) {
		for {
			for (len(controller.sentMessages) != 0) && (time.Now().Sub(controller.sentMessages[0]) > time.Second) {
				controller.sentMessages = controller.sentMessages[1:]
			}

			if len(controller.sentMessages) < maxMessagesCountPerSecond {
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
