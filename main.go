package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/valeyard1/lastfmanonbot/lastfmanonbot"
)

var (
	token = os.Getenv("TELEGRAM_BOT_TOKEN")
)

func main() {

	if token == "" {
		log.Fatal("No token has been provided for bot to work. Provide the TELEGRAM_BOT_TOKEN environment variable")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Initialize LastFM Api
		lastfmanonbot.CreateLastfmApi()
		switch update.Message.Command() {
		case "help":
			msg.Text = lastfmanonbot.HelpMessage()
		case "status":
			msg.Text = lastfmanonbot.GetNowPlaying(update.Message.CommandArguments())
		default:
			msg.Text = lastfmanonbot.HelpMessage()
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
