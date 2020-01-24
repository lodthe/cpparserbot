// Package logger implements different logging agents
package logger

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lodthe/cpparserbot/controller"
	"github.com/lodthe/cpparserbot/helper"
)

// TelegramLoggers sends messages to the telegram channel
// specified by channel id
type TelegramLogger struct {
	ChannelID  int64
	Controller *controller.TelegramController
}

// Info sends messages with informative context
func (logger *TelegramLogger) Info(text string) {
	msg := tgbotapi.NewMessage(logger.ChannelID, fmt.Sprintf("*[info]*: %s", text))
	logger.Controller.Send(helper.PrepareMessageConfig(&msg))
}

// Error sends messages about errors
func (logger *TelegramLogger) Error(text string) {
	msg := tgbotapi.NewMessage(logger.ChannelID, fmt.Sprintf("*[error]*: %s", text))
	logger.Controller.Send(helper.PrepareMessageConfig(&msg))
}
