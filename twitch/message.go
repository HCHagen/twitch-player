package twitch

import (
	"time"
)

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
