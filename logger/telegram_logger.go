// Package logger implements different logging agents
package logger

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lodthe/cpparserbot/label"

	"github.com/lodthe/cpparserbot/controller"
	"github.com/lodthe/cpparserbot/helper"
)

// TelegramLoggers sends messages to the telegram channel
// specified by channel id
type TelegramLogger struct {
	ChannelID  int64
	Controller controller.Controller
}

// Info sends messages with informative context
func (logger *TelegramLogger) Info(text string) {
	msg := tgbotapi.NewMessage(logger.ChannelID, fmt.Sprintf(label.InfoFormat, text))
	helper.PrepareMessage(&msg)
	logger.Controller.Send(&msg)
}

// Error sends messages about errors
func (logger *TelegramLogger) Error(text string) {
	msg := tgbotapi.NewMessage(logger.ChannelID, fmt.Sprintf(label.ErrorFormat, text))
	helper.PrepareMessage(&msg)
	logger.Controller.Send(&msg)
}
