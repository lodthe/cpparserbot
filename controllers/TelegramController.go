package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//TelegramController controls request sending to
//prevent against exceeding Telegram limits
type TelegramController struct {
	Bot *tgbotapi.BotAPI
}

//Send sends something to chat
func (controller *TelegramController) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return controller.Bot.Send(c)
}
