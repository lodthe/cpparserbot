package helpers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

//PrepareMessageConfig disables web page preview
//and sets parse mode to Markdown
func PrepareMessageConfig(config tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	config.ParseMode = tgbotapi.ModeMarkdown
	config.DisableWebPagePreview = true
	return config
}
