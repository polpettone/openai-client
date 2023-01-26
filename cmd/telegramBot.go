package cmd

import (
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot() error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("CLEVER_MAN_TELEGRAM_BOT"))

	if err != nil {
		return err
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	provider := Provider{}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		fmt.Printf("MessageFrom: %s\n", update.Message.From)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		triggerWord := "openAI"
		prompt := ""
		response := "no prompt to openai"

		if strings.Contains(msg.Text, triggerWord) {

			prompt = strings.Replace(msg.Text, triggerWord, "", -1)

			response, err = provider.Prompt(prompt, "text-davinci-003", 0.7, 256)
			if err != nil {
				return err
			}
		}

		msg.ReplyToMessageID = update.Message.MessageID
		msg.Text = response

		if _, err := bot.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
