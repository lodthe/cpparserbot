package helpers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//PrepareMessageConfigForGroup disables web page preview,
//sets parse mode to Markdown and removes keyboards
//in group chats
func PrepareMessageConfig(config *tgbotapi.MessageConfig) *tgbotapi.MessageConfig {
	config.ParseMode = tgbotapi.ModeMarkdown
	config.DisableWebPagePreview = true

	if config.ChatID < 0 {
		config.ReplyMarkup = nil
	}

	return config
}

//PreparePhotoConfig disables web page preview,
//sets parse mode to Markdown and removes keyboards
//in group chats
func PreparePhotoConfig(config *tgbotapi.PhotoConfig) *tgbotapi.PhotoConfig {
	config.ParseMode = tgbotapi.ModeMarkdown

	if config.ChatID < 0 {
		config.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	}

	return config
}

//GetChatID tries to get chat ID from update object
func GetChatID(update *tgbotapi.Update) int64 {
	switch true {
	case update.Message != nil:
		return update.Message.Chat.ID
	case update.CallbackQuery != nil:
		return update.CallbackQuery.Message.Chat.ID
	default:
		return 0
	}
}

//GetTelegramProfileURL parses update owner and for user chats
//create URL to their profiles
//P.S. Chat ID > 0 for user chats
func GetTelegramProfileURL(update *tgbotapi.Update) string {
	ID := GetChatID(update)
	if ID < 0 {
		return fmt.Sprintf("%d", ID)
	} else {
		return fmt.Sprintf("[%v](tg://user?id=%v)", ID, ID)
	}

}
