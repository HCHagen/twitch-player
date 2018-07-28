package twitch

import (
	"fmt"
	"time"
)

type GameInfo struct {
	Game     Game   `json:"game"`
	Viewers  uint64 `json:viewers"`
	Channels uint64 `json:channels"`
}

type GameListResult struct {
	Total uint64     `json:"_total"`
	Games []GameInfo `json:"top"`
}

type Featured struct {
	Image     string `json:"image"`
	Priority  uint64 `json:"priority"`
	Scheduled bool   `json:"scheduled"`
	Sponsored bool   `json:"sponsored"`
	Stream    Stream `json:"stream"`
	Text      string `json:"text"`
	Title     string `json:"title"`
}

type FeaturedListResult struct {
	Featured []Featured `json:"featured"`
}

type StreamListResult struct {
	Total   uint64   `json:"_total"`
	Streams []Stream `json:"streams"`
}

type ErrorResponse struct {
	error
	StatusCode uint64 `json:"status"`
	Status     string `json:"error"`
	Message    string `json:"message"`
}

type ChannelSearchResult struct {
	Total    uint64    `json:"_total"`
	Channels []Channel `json:"channels"`
}

type GameSearchResult struct {
	Games []Game `json:"games"`
}

type StreamSearchResult StreamListResult

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("HTTP %d - %s: %s", e.StatusCode, e.Status, e.Message)
}

type Stream struct {
	Id          uint64    `json:"_id"`
	Game        string    `json:"game"`
	Viewers     uint64    `json:"viewers"`
	VideoHeight uint64    `json:"video_height"`
	AvgFps      float64   `json:"average_fps"`
	Delay       uint64    `json:"delay"`
	Created     time.Time `json:"created_at"`
	IsPlaylist  bool      `json:"is_playlist"`
	Channel     Channel   `json:"channel"`
}

type Game struct {
	Name          string `json:"name"`
	Popularity    uint64 `json:"popularity"`
	Id            uint64 `json:"_id"`
	GiantbombId   uint64 `json:"giantbomb_id"`
	LocalizedName string `json:"localized_name"`
	Locale        string `json:"locale"`
}

type Channel struct {
	Mature              bool      `json:"mature"`
	Status              string    `json:"status"`
	BroadcasterLanguage string    `json:"broadcaster_language"`
	DisplayName         string    `json:"display_name"`
	Game                string    `json:"game"`
	Language            string    `json:"language"`
	Id                  uint64    `json:"_id"`
	Name                string    `json:"name"`
	Created             time.Time `json:"created_at"`
	Updated             time.Time `json:"updated_at"`
	Partner             bool      `json:"partner"`
	Logo                string    `json:"logo"`
	VideoBanner         string    `json:"video_banner"`
	ProfileBanner       string    `json:"profile_banner"`
	Url                 string    `json:"url"`
	Views               uint64    `json:"views"`
	Followers           uint64    `json:"followers"`

	BroadcasterType string `json:"broadcaster_type,omitempty"`
	Description     string `json:"description,omitempty"`
	PrivateVideo    bool   `json:"private_video,omitempty"`
	PrivacyOptions  bool   `json:"privacy_options_enabled,omitempty"`
}

type Links struct {
	Self    string `json:"self"`
	Channel string `json:"channel"`
}

type StreamData struct {
	Stream Stream `json:"stream"`
	Links  Links  `json:"_links"`
}

type Token struct {
	Token     string `json:"token"`
	Signature string `json:"sig"`
}

type StreamUrl struct {
	Bandwidth  uint32
	Quality    string
	Resolution string
	URI        string
}
