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

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
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
			response, err = Questioner(prompt, "text-davinci-003", 0.7, 256)
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
