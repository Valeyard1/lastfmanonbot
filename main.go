package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/valeyard1/lastfmanonbot/lastfmanonbot"
	tele "gopkg.in/telebot.v3"
)

var (
	token = os.Getenv("TELEGRAM_BOT_TOKEN")
)

func main() {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	lastfmanonbot.CreateLastfmApi()
	b.Handle(tele.OnQuery, func(c tele.Context) error {
		user := c.Query().Text

		urls := []string{
			lastfmanonbot.GetNowPlayingAlbumArt(user),
		}
		song := lastfmanonbot.GetNowPlayingSong(user)
		artist := lastfmanonbot.GetNowPlayingArtist(user)
		album := lastfmanonbot.GetNowPlayingAlbum(user)
		albumURL := lastfmanonbot.GetNowPlayingAlbumURL(user)
		verbalTense := lastfmanonbot.GetNowPlayingVerbalTense(user)
		tags := lastfmanonbot.GetNowPlayingSongTags(user)

		nowPlayingTrack := fmt.Sprintf("I%s listening to\n🎧 %s - %s [%s]\n\n%s",
			verbalTense, artist, song, album, tags)

		results := make(tele.Results, len(urls)) // []tele.Result
		for i, url := range urls {
			result := &tele.ArticleResult{
				Title:       song,
				Text:        nowPlayingTrack,
				Description: artist,
				URL:         albumURL,
				ThumbURL:    url,
				HideURL:     true,
			}
			results[i] = result
			// needed to set a unique string ID for each result
			results[i].SetResultID(strconv.Itoa(i))
		}

		return c.Answer(&tele.QueryResponse{
			Results:   results,
			CacheTime: 60, // a minute
		})

		// no log since it's anonymous
	})

	b.Start()
}
