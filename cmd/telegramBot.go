package cmd

import (
	"encoding/json"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/polpettone/openai-client/cmd/config"
)

func StartBot(telegramBotToken, contextMemoryID string) error {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)

	if err != nil {
		return err
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	provider := NewProvider(3000, true, contextMemoryID)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		config.Logger.
			Info().
			Str("MessageFrom", update.Message.From.String()).
			Send()

		messageJson, err := json.Marshal(update.Message)

		if err != nil {
			return err
		}

		config.Logger.
			Info().
			RawJSON("Message", messageJson).
			Send()

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		triggerClearContext := "clearContext"
		triggerWord := "openAI"
		prompt := ""
		response := "no prompt to openai"

		if strings.Contains(msg.Text, triggerWord) {

			if strings.Contains(msg.Text, triggerClearContext) {
				err := provider.ClearContext()
				if err != nil {
					return err
				}
				response = "context cleared"
			} else {
				prompt = strings.Replace(msg.Text, triggerWord, "", -1)
				response, err = provider.Prompt(prompt, "text-davinci-003", 0.7, 3000)
				if err != nil {
					return err
				}
			}
		}

		msg.ReplyToMessageID = update.Message.MessageID
		if response == "" {
			response = "no response"
		}
		msg.Text = response

		if _, err := bot.Send(msg); err != nil {
			return err
		}

	}
	return nil
}
