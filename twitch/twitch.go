package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/grafov/m3u8"
)

const (
	AccessTokenUrl     = "https://api.twitch.tv/api/channels/%s/access_token"
	GetStreamUrl       = "https://api.twitch.tv/kraken/streams/%s"
	StreamGeneratorUrl = "https://usher.ttvnw.net/api/channel/hls/%s.m3u8?player=twitchweb&token=%s&sig=%s&allow_audio_only=true&allow_source=true&type=any&allow_spectre=false&p=%d"
)

type Client interface {
	GetStreamData(channel string) (StreamData, error)
	GetStreamUrls(channel string) ([]StreamUrl, error)
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
		msg, _ := ioutil.ReadAll(res.Body)
		return tok, fmt.Errorf("Getting token: HTTP %s: %s", res.Status, msg)
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
		msg, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("Getting streams: HTTP %s: %s", res.Status, msg)
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

func (c *twitchClient) GetStreamData(channel string) (sd StreamData, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(GetStreamUrl, channel), nil)
	if err != nil {
		return sd, err
	}
	req.Header.Add("Client-ID", c.clientId)

	res, err := c.Do(req)
	if err != nil {
		return sd, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(res.Body)
		return sd, fmt.Errorf("Getting stream data: HTTP %s: %s", res.Status, msg)
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
