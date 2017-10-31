package webapi

import (
	"github.com/r3boot/go-musicbot/lib/config"
	"github.com/r3boot/go-musicbot/lib/mp3lib"
	"github.com/r3boot/go-musicbot/lib/mpdclient"
	"github.com/r3boot/go-musicbot/lib/ytclient"
)

const MAX_PLAYLIST_LENGTH int = 8192

type TemplateData struct {
	Title  string
	Stream string
}

type WebApi struct {
	config *config.MusicBotConfig
	mpd    *mpdclient.MPDClient
	mp3    *mp3lib.MP3Library
	yt     *youtubeclient.YoutubeClient
}

type ClientRequest struct {
	Operation string
}

type SearchRequest struct {
	Operation string
	Query string
}

type NowPlaying struct {
	Title    string
	Duration string
	Rating   int
}

type NowPlayingResp struct {
	Data NowPlaying
	Pkt  string
}

type GetFilesResp struct {
	Data []string
	Pkt string
}