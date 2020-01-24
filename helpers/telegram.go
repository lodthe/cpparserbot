package helpers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

//PrepareMessageConfigForGroup disables web page preview,
//sets parse mode to Markdown and removes keyboards
//in group chats
func PrepareMessageConfig(config *tgbotapi.MessageConfig) *tgbotapi.MessageConfig {
	config.ParseMode = tgbotapi.ModeMarkdown
	config.DisableWebPagePreview = true

	if config.ChatID < 0 {
		config.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
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
