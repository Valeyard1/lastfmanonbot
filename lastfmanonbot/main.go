package lastfmanonbot

import (
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

func GetNowPlayingSong(user string) string {
	userMap := lastfm.P{
		"user": user,
	}

	nowPlaying, _ := api.User.GetRecentTracks(userMap)
	for _, v := range nowPlaying.Tracks {
		return v.Name
	}
	return ""
}

func GetNowPlayingArtist(user string) string {
	userMap := lastfm.P{
		"user": user,
	}

	nowPlaying, _ := api.User.GetRecentTracks(userMap)
	for _, v := range nowPlaying.Tracks {
		return v.Artist.Name
	}
	return ""
}

func GetNowPlayingAlbum(user string) string {
	userMap := lastfm.P{
		"user": user,
	}

	nowPlaying, _ := api.User.GetRecentTracks(userMap)
	for _, v := range nowPlaying.Tracks {
		return v.Album.Name
	}
	return ""
}

func GetNowPlayingAlbumURL(user string) string {
	albumMap := lastfm.P{
		"mbid": getNowPlayingAlbumMbid(user),
	}

	albumInfo, _ := api.Album.GetInfo(albumMap)
	albumArt := albumInfo.Url

	return albumArt
}

func GetNowPlayingAlbumArt(user string) string {
	// TODO: get album art by mbid and return nil if no album art
	// mbid := getNowPlayingAlbumMbid(user)

	albumMap := lastfm.P{
		"track":  GetNowPlayingSong(user),
		"artist": GetNowPlayingArtist(user),
	}

	trackInfo, _ := api.Track.GetInfo(albumMap)
	albumArt := trackInfo.Album

	for _, images := range albumArt.Images {
		if images.Size == "extralarge" {
			return images.Url
		}
	}

	return "none"
}

func getNowPlayingAlbumMbid(user string) string {
	userMap := lastfm.P{
		"user": user,
	}

	nowPlaying, _ := api.User.GetRecentTracks(userMap)
	for _, v := range nowPlaying.Tracks {
		return v.Album.Mbid
	}
	return ""
}

func GetNowPlayingVerbalTense(user string) string {
	userMap := lastfm.P{
		"user": user,
	}

	nowPlaying, _ := api.User.GetRecentTracks(userMap)
	for _, v := range nowPlaying.Tracks {
		if v.NowPlaying == "true" {
			return "'m"
		} else {
			return " was"
		}
	}
	return ""
}

func HelpMessage() string {
	message := `Available commands:
/status [username] - Your rencently played song
/help - Display this message
	`
	return message
}
