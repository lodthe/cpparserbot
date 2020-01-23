package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/lodthe/cpparserbot/controllers"
	"github.com/lodthe/cpparserbot/handlers"
	"github.com/lodthe/cpparserbot/loggers"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	controller := &controllers.TelegramController{Bot: bot}
	controller.Init()
	controller.Run()

	channelID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHANNEL_CHAT_ID"), 10, 64)
	logger := &loggers.TelegramLogger{
		ChannelID:  channelID,
		Controller: controller,
	}

	logger.Info("Bot was started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handlers.DispatchMessage(update, controller, logger)
		}

		if update.CallbackQuery != nil {
			handlers.DispatchCallback(update, controller, logger)
		}
	}
}
