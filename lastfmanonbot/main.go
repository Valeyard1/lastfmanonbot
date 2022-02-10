package lastfmanonbot

import (
	"fmt"
	"log"
	"os"

	"github.com/shkh/lastfm-go/lastfm"
)

var api *lastfm.Api

func CreateLastfmApi() bool {
	apiKey := os.Getenv("LASTFM_APIKEY")
	apiSharedSecret := os.Getenv("LASTFM_SHAREDSECRET")
	if apiKey == "" || apiSharedSecret == "" {
		log.Fatal(`No API keys has been provided for bot to integrate with LastFM.
		Provide the LASTFM_APIKEY and LASTFM_SHAREDSECRET environment variable`)
	}

	api = lastfm.New(apiKey, apiSharedSecret)
	if api == nil {
		log.Panic(api)
	}

	return true
}

func GetNowPlaying(user string) string {
	var nowPlayingTrack string

	userMap := lastfm.P{
		"user": user,
	}

	nowPlaying, _ := api.User.GetRecentTracks(userMap)
	for _, v := range nowPlaying.Tracks {
		var verbalTense string
		if v.NowPlaying == "true" {
			verbalTense = fmt.Sprintf("'m")
		} else {
			verbalTense = fmt.Sprintf(" was")
		}

		nowPlayingTrack = fmt.Sprintf("I%s listening to:\n%s - %s [%s]",
			verbalTense, v.Artist.Name, v.Name, v.Album.Name)
		break
	}

	return nowPlayingTrack
}

func HelpMessage() string {
	message := `Available commands:
/status [username] - Your rencently played song
/help - Display this message
	`
	return message
}
