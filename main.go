package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shkh/lastfm-go/lastfm"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
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

		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "status":
			msg.Text = getNowPlaying(update.Message.CommandArguments())
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func getNowPlaying(user string) string {
	apiKey := os.Getenv("LASTFM_APIKEY")
	apiSharedSecret := os.Getenv("LASTFM_SHAREDSECRET")
	var nowPlayingTrack string

	api := lastfm.New(apiKey, apiSharedSecret)
	userMap := lastfm.P{
		"user": user,
	}

	nowPlaying, _ := api.User.GetRecentTracks(userMap)
	for _, v := range nowPlaying.Tracks {
		nowPlayingTrack = fmt.Sprintf("I'm listening to:\n%s - %s [%s]", v.Artist.Name, v.Name, v.Album.Name)
		break
	}

	return nowPlayingTrack
}
