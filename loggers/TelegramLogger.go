package loggers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lodthe/cpparserbot/controllers"
	"github.com/lodthe/cpparserbot/helpers"
)

//TelegramLoggers sends messages to the telegram channel
//specified by channel id
type TelegramLogger struct {
	ChannelID  int64
	Controller *controllers.TelegramController
}

//Info sends messages with informative context
func (logger *TelegramLogger) Info(text string) {
	logger.Controller.Send(
		helpers.PrepareMessageConfig(
			tgbotapi.NewMessage(
				logger.ChannelID,
				fmt.Sprintf("*[info]*: %s", text),
			)))
}

//Error sends messages about errors
func (logger *TelegramLogger) Error(text string) {
	logger.Controller.Send(
		helpers.PrepareMessageConfig(
			tgbotapi.NewMessage(
				logger.ChannelID,
				fmt.Sprintf("*[error]*: %s", text),
			)))
}
