package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lodthe/cpparserbot/controllers"
	"github.com/lodthe/cpparserbot/labels"
	"github.com/lodthe/cpparserbot/loggers"
)

//DispatchMessage identifies type of message and sends response based on this type
func DispatchMessage(update tgbotapi.Update, controller *controllers.TelegramController, logger *loggers.TelegramLogger) {
	switch update.Message.Text {
	case "/start":
		logger.Info(fmt.Sprintf("User %v sent /start", update.Message.Chat.ID))
		controller.Send(handleStart(update))
	}
}

//handleStart handles /start command
func handleStart(update tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, labels.START)
}
