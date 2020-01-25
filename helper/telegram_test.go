package helper

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestGetChatID(t *testing.T) {
	in := []tgbotapi.Update{
		{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
			},
		},

		{
			CallbackQuery: &tgbotapi.CallbackQuery{
				Message: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: 2,
					},
				},
			},
		},

		{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 3,
				},
			},
			CallbackQuery: &tgbotapi.CallbackQuery{
				Message: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: 3,
					},
				},
			},
		},

		{},
	}
	want := []int64{1, 2, 3, 0}

	for i := range in {
		ID := GetChatID(&in[i])
		if ID != want[i] {
			t.Errorf("For %v expected %d, got %d", in[i], want[i], ID)
		}
	}
}

func TestPrepareMessage(t *testing.T) {
	in := []tgbotapi.Chattable{
		&tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      1000,
				ReplyMarkup: tgbotapi.ReplyKeyboardMarkup{},
			},
			Text:                  "First",
			ParseMode:             tgbotapi.ModeHTML,
			DisableWebPagePreview: true,
		},

		&tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      -1000,
				ReplyMarkup: tgbotapi.ReplyKeyboardMarkup{},
			},
			Text:                  "Second",
			DisableWebPagePreview: true,
		},

		&tgbotapi.PhotoConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      1000,
					ReplyMarkup: tgbotapi.ReplyKeyboardMarkup{},
				},
			},
			Caption: "Third",
		},

		&tgbotapi.PhotoConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      -1000,
					ReplyMarkup: tgbotapi.ReplyKeyboardMarkup{},
				},
			},
			Caption: "Third",
		},

		&tgbotapi.DocumentConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      1000,
					ReplyMarkup: tgbotapi.ReplyKeyboardMarkup{},
				},
			},
			Caption:   "Fourth",
			ParseMode: tgbotapi.ModeHTML,
		},

		&tgbotapi.DocumentConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      -1000,
					ReplyMarkup: tgbotapi.ReplyKeyboardMarkup{},
				},
			},
			Caption:   "Fourth",
			ParseMode: tgbotapi.ModeHTML,
		},
	}

	for _, c := range in {
		PrepareMessage(c)

		switch c.(type) {
		case *tgbotapi.MessageConfig:
			v := c.(*tgbotapi.MessageConfig)
			if (v.DisableWebPagePreview == false) || (v.ParseMode != tgbotapi.ModeMarkdown) ||
				((v.ReplyMarkup != nil) && (v.ChatID < 0)) {
				t.Errorf("Expected DisableWebPagePreview == true, ParseMode == Markdown, ReplyMarkup = nil (if id < 0), got %v", v)
			}

		case *tgbotapi.PhotoConfig:
			v := c.(*tgbotapi.PhotoConfig)
			if (v.ParseMode != tgbotapi.ModeMarkdown) || ((v.ReplyMarkup != nil) && (v.ChatID < 0)) {
				t.Errorf("Expected ParseMode == Markdown, ReplyMarkup = nil (if id < 0), got %v", v)
			}

		case *tgbotapi.DocumentConfig:
			v := c.(*tgbotapi.DocumentConfig)
			if (v.ParseMode != tgbotapi.ModeMarkdown) || ((v.ReplyMarkup != nil) && (v.ChatID < 0)) {
				t.Errorf("Expected ParseMode == Markdown, ReplyMarkup = nil (if id < 0), got %v", v)
			}
		}
	}
}
