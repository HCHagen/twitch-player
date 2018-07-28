package twitch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/grafov/m3u8"
)

const (
	AccessTokenUrl = "https://api.twitch.tv/api/channels/%s/access_token"

	GetStreamUrl = "https://api.twitch.tv/kraken/streams/%d"

	SearchChannelUrl = "https://api.twitch.tv/kraken/search/channels?query=%s&limit=%d"
	SearchGameUrl    = "https://api.twitch.tv/kraken/search/games?query=%s&type=suggest"
	SearchStreamUrl  = "https://api.twitch.tv/kraken/search/streams?query=%s&limit=%d"

	ListGamesUrl       = "https://api.twitch.tv/kraken/games/top?limit=%d"
	ListFeaturedUrl    = "https://api.twitch.tv/kraken/streams/featured?limit=%d"
	ListGameStreamsUrl = "https://api.twitch.tv/kraken/streams/?game=%s&limit=%d"

	StreamGeneratorUrl = "https://usher.ttvnw.net/api/channel/hls/%s.m3u8?player=twitchweb&token=%s&sig=%s&allow_audio_only=true&allow_source=true&type=any&allow_spectre=false&p=%d"

	KrakenApiAcceptHeader = "application/vnd.twitchtv.v5+json"
)

type Client interface {
	GetStreamData(channelId uint64) (StreamData, error)
	GetStreamUrls(channel string) ([]StreamUrl, error)

	GetChannelSearch(channel string, num int) (ChannelSearchResult, error)
	GetGameSearch(game string) (GameSearchResult, error)
	GetStreamSearch(channel string, num int) (StreamSearchResult, error)

	GetGameList(num int) (GameListResult, error)
	GetFeaturedList(num int) (FeaturedListResult, error)
	GetStreamList(game string, num int) (StreamListResult, error)
}

type twitchClient struct {
	*http.Client

	clientId string
}

func NewTwitchClient(clientId string, httpTimeout time.Duration) (Client, error) {
	httpClient, err := newHttp2Client(httpTimeout, nil)
	if err != nil {
		return nil, err
	}
	return &twitchClient{
		Client:   httpClient,
		clientId: clientId,
	}, nil
}

func (c *twitchClient) getStreamToken(channel string) (tok Token, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(AccessTokenUrl, channel), nil)
	if err != nil {
		return tok, err
	}
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return tok, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return tok, ErrUnmarshal("Getting token", res)
	}

	return tok, json.NewDecoder(res.Body).Decode(&tok)
}

func (c *twitchClient) getStreamUrls(channel string, tok Token) (*m3u8.MasterPlaylist, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(StreamGeneratorUrl, channel, url.QueryEscape(tok.Token), tok.Signature, time.Now().UnixNano()), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, ErrUnmarshal("Getting streams", res)
	}

	p, listType, err := m3u8.DecodeFrom(res.Body, false)
	if err != nil {
		return nil, err
	}

	if listType != m3u8.MASTER {
		return nil, fmt.Errorf("Stream offline or does not exist")
	}

	return p.(*m3u8.MasterPlaylist), nil
}

func (c *twitchClient) GetChannelSearch(channel string, num int) (sr ChannelSearchResult, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(SearchChannelUrl, url.QueryEscape(channel), num), nil)
	if err != nil {
		return sr, err
	}
	req.Header.Add("Accept", KrakenApiAcceptHeader)
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return sr, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return sr, ErrUnmarshal("Searching for channel", res)
	}

	return sr, json.NewDecoder(res.Body).Decode(&sr)
}

func (c *twitchClient) GetGameSearch(game string) (sr GameSearchResult, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(SearchGameUrl, url.QueryEscape(game)), nil)
	if err != nil {
		return sr, err
	}
	req.Header.Add("Accept", KrakenApiAcceptHeader)
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return sr, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return sr, ErrUnmarshal("Searching for game", res)
	}

	return sr, json.NewDecoder(res.Body).Decode(&sr)
}

func (c *twitchClient) GetStreamSearch(channel string, num int) (sr StreamSearchResult, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(SearchStreamUrl, url.QueryEscape(channel), num), nil)
	if err != nil {
		return sr, err
	}
	req.Header.Add("Accept", KrakenApiAcceptHeader)
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return sr, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return sr, ErrUnmarshal("Searching for stream", res)
	}

	return sr, json.NewDecoder(res.Body).Decode(&sr)
}

func (c *twitchClient) GetGameList(num int) (lr GameListResult, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(ListGamesUrl, num), nil)
	if err != nil {
		return lr, err
	}
	req.Header.Add("Accept", KrakenApiAcceptHeader)
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return lr, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return lr, ErrUnmarshal("Listing games", res)
	}

	return lr, json.NewDecoder(res.Body).Decode(&lr)
}

func (c *twitchClient) GetFeaturedList(num int) (lr FeaturedListResult, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(ListFeaturedUrl, num), nil)
	if err != nil {
		return lr, err
	}
	req.Header.Add("Accept", KrakenApiAcceptHeader)
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return lr, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return lr, ErrUnmarshal("Listing featured streams", res)
	}

	return lr, json.NewDecoder(res.Body).Decode(&lr)
}

func (c *twitchClient) GetStreamList(game string, num int) (lr StreamListResult, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(ListGameStreamsUrl, url.QueryEscape(game), num), nil)
	if err != nil {
		return lr, err
	}
	req.Header.Add("Accept", KrakenApiAcceptHeader)
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return lr, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return lr, ErrUnmarshal("Listing streams", res)
	}

	return lr, json.NewDecoder(res.Body).Decode(&lr)
}

func (c *twitchClient) GetStreamData(channelId uint64) (sd StreamData, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(GetStreamUrl, channelId), nil)
	if err != nil {
		return sd, err
	}
	req.Header.Add("Accept", KrakenApiAcceptHeader)
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return sd, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return sd, ErrUnmarshal("Getting stream data", res)
	}

	return sd, json.NewDecoder(res.Body).Decode(&sd)
}

func (c *twitchClient) GetStreamUrls(channel string) (streams []StreamUrl, err error) {
	tok, err := c.getStreamToken(channel)
	if err != nil {
		return streams, err
	}
	pl, err := c.getStreamUrls(channel, tok)
	if err != nil {
		return streams, err
	}

	for _, variant := range pl.Variants {
		streams = append(streams, StreamUrl{
			Bandwidth:  variant.VariantParams.Bandwidth,
			Quality:    variant.VariantParams.Video,
			Resolution: variant.VariantParams.Resolution,
			URI:        variant.URI,
		})
	}

	return streams, nil
}

func ErrUnmarshal(action string, res *http.Response) error {
	var err ErrorResponse
	if err := json.NewDecoder(res.Body).Decode(&err); err != nil {
		return fmt.Errorf("%s: HTTP %s", action, res.Status)
	}
	return fmt.Errorf("%s: %s", action, err.Error())
}
