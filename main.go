// Package main runs bot and additional things like logger, controller
package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lodthe/cpparserbot/api"
	"github.com/lodthe/cpparserbot/controller"
	"github.com/lodthe/cpparserbot/handlers"
	"github.com/lodthe/cpparserbot/logger"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	controller := &controller.TelegramController{Bot: bot}
	controller.Init()
	controller.Run()

	channelID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHANNEL_CHAT_ID"), 10, 64)
	logger := &logger.TelegramLogger{
		ChannelID:  channelID,
		Controller: controller,
	}

	binanceAPI := &api.Binance{}
	binanceAPI.Init(os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_SECRET_KEY"))

	messageDispatcher := handlers.NewMessageDispatcher(controller, logger, binanceAPI)

	logger.Info("Bot was started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go messageDispatcher.Dispatch(&update)
		}
	}
}
