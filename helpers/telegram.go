package helpers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

//PrepareMessageConfig disables web page preview
//and sets parse mode to Markdown
func PrepareMessageConfig(config tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	config.ParseMode = tgbotapi.ModeMarkdown
	config.DisableWebPagePreview = true
	return config
}

//PrepareMessageConfigForGroup disables web page preview,
//sets parse mode to Markdown and removes keyboards
//in group chats
func PrepareMessageConfigForGroup(config *tgbotapi.MessageConfig, update *tgbotapi.Update) *tgbotapi.MessageConfig {
	config.ParseMode = tgbotapi.ModeMarkdown
	config.DisableWebPagePreview = true

	if update.Message.Chat.ID < 0 {
		config.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	}

	return config
}
