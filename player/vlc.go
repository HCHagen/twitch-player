package player

import (
	"fmt"

	vlc "github.com/adrg/libvlc-go"
)

type vlcPlayer struct {
	player *vlc.Player

	loadedMedia *vlc.Media
}

func NewVlcPlayer() (Player, error) {
	if err := vlc.Init("--quiet"); err != nil {
		return nil, err
	}

	player, err := vlc.NewPlayer()
	if err != nil {
		vlc.Release()
		return nil, err
	}

	return &vlcPlayer{
		player: player,
	}, nil
}

func (p *vlcPlayer) Reset() (res error) {
	if p.player.IsPlaying() {
		if err := p.player.Stop(); err != nil {
			res = err
		}
	}
	if p.loadedMedia != nil {
		if err := p.loadedMedia.Release(); err != nil {
			res = err
		}
	}

	return res
}

func (p *vlcPlayer) Close() (res error) {
	res = p.Reset()

	if err := p.player.Release(); err != nil {
		res = err
	}
	if err := vlc.Release(); err != nil {
		res = err
	}

	return res
}

func (p *vlcPlayer) LoadFromUrl(url string) (err error) {
	p.Reset()

	p.loadedMedia, err = p.player.LoadMediaFromURL(url)
	if err != nil {
		return err
	}

	return nil
}

func (p *vlcPlayer) LoadFromFile(file string) (err error) {
	p.Reset()

	p.loadedMedia, err = p.player.LoadMediaFromPath(file)
	if err != nil {
		return err
	}

	return nil
}

func (p *vlcPlayer) Play() error {
	if p.loadedMedia == nil {
		return fmt.Errorf("Cannot play: No media loaded")
	}

	return p.player.Play()
}

func (p *vlcPlayer) Stop() error {
	if !p.player.IsPlaying() {
		return fmt.Errorf("Cannot stop: No media playing")
	}

	return p.player.Stop()
}

func (p *vlcPlayer) EnterFullscreen() error {
	if fs, err := p.player.IsFullScreen(); err != nil {
		return fmt.Errorf("Cannot enter fullscreen: %s", err.Error())
	} else if fs {
		return fmt.Errorf("Cannot enter fullscreen: Is already fullscreen")
	}

	return p.player.ToggleFullScreen()
}

func (p *vlcPlayer) ExitFullscreen() error {
	if fs, err := p.player.IsFullScreen(); err != nil {
		return fmt.Errorf("Cannot exit fullscreen: %s", err.Error())
	} else if !fs {
		return fmt.Errorf("Cannot exit fullscreen: No fullscreen to exit")
	}

	return p.player.ToggleFullScreen()
}
