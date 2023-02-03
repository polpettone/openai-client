package commands

import (
	"encoding/json"
	"strings"

	provider2 "github.com/polpettone/openai-client/cmd/provider"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/polpettone/openai-client/cmd/logging"
)

func StartBot(telegramBotToken, contextMemoryID string) error {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)

	if err != nil {
		return err
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	provider := provider2.NewProvider(3000, true, contextMemoryID)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		logging.Logger.
			Info().
			Str("MessageFrom", update.Message.From.String()).
			Send()

		err := logMessage(update.Message)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		triggerClearContext := "clearContext"
		triggerWord := "openAI"
		prompt := ""
		response := "no prompt to openai"
		shutdownTrigger := "machAusDitDing"

		if strings.Contains(msg.Text, shutdownTrigger) {
			return nil
		}

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

		message, err := bot.Send(msg)

		if err != nil {
			return err
		}

		err = logMessage(&message)
		if err != nil {
			return err
		}

	}
	return nil
}

func logMessage(msg *tgbotapi.Message) error {
	messageJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	logging.Logger.
		Info().
		RawJSON("Message", messageJson).
		Send()
	return nil
}
