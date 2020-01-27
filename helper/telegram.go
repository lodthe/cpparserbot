package helper

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// PrepareMessage prepares message for sending:
//  disables web page preview,
//  sets parse mode to Markdown,
//  removes keyboard in group chats
func PrepareMessage(config tgbotapi.Chattable) tgbotapi.Chattable {
	switch config.(type) {
	case *tgbotapi.MessageConfig:
		message := config.(*tgbotapi.MessageConfig)
		message.ParseMode = tgbotapi.ModeMarkdown
		message.DisableWebPagePreview = true
		if message.ChatID < 0 {
			message.ReplyMarkup = nil
		}

	case *tgbotapi.PhotoConfig:
		photo := config.(*tgbotapi.PhotoConfig)
		photo.ParseMode = tgbotapi.ModeMarkdown
		if photo.ChatID < 0 {
			photo.ReplyMarkup = nil
		}

	case *tgbotapi.DocumentConfig:
		document := config.(*tgbotapi.DocumentConfig)
		document.ParseMode = tgbotapi.ModeMarkdown
		if document.ChatID < 0 {
			document.ReplyMarkup = nil
		}
	default:
		log.Fatalf("Cannot recognize Config type: %T", config)
	}

	return config
}

// GetChatID tries to get chat ID from update object
func GetChatID(update *tgbotapi.Update) int64 {
	switch {
	case update.Message != nil:
		return update.Message.Chat.ID
	case update.CallbackQuery != nil:
		return update.CallbackQuery.Message.Chat.ID
	default:
		return 0
	}
}

// GetTelegramProfileURL parses update owner and for user chats
// create URL to their profiles
// P.S. Chat ID > 0 for user chats
func GetTelegramProfileURL(update *tgbotapi.Update) string {
	ID := GetChatID(update)
	if ID < 0 {
		return fmt.Sprintf("%d", ID)
	} else {
		return fmt.Sprintf("[%v](tg://user?id=%v)", ID, ID)
	}
}
