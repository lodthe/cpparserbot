package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lodthe/cpparserbot/controllers"
	"github.com/lodthe/cpparserbot/loggers"
)

//DispatchCallback identifies type of callback and sends response based on the callback type
func DispatchCallback(update tgbotapi.Update, controller *controllers.TelegramController, logger *loggers.TelegramLogger) {

}
