package helpers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

//PrepareMessageConfigForSending disables web page preview
//and sets parse mode to Markdown
func PrepareMessageConfigForSending(config tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	config.ParseMode = tgbotapi.ModeMarkdown
	config.DisableWebPagePreview = true
	return config
}
