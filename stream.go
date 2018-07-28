package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"
)

func onStream(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("Please provide a channel name")
	}

	channelName := ctx.Args()[0]

	sr, err := twitchClient().GetChannelSearch(channelName, DefaultChannelResultLen)
	if err != nil {
		return err
	}

	if len(sr.Channels) == 0 {
		return fmt.Errorf("No channels found for %s", channelName)
	}

	var chanSelection int
	if sr.Channels[0].Name != channelName {
		fmt.Printf("No channels found for %s. Did you possibly mean...\n\n", ctx.Args()[0])
		for i, c := range sr.Channels {
			fmt.Printf("[%d - %s] %s (id %d)?\n", i, c.Name, c.DisplayName, c.Id)
		}
		chanSelection = getNumericInput(fmt.Sprintf("\nSelect stream channel: [0-%d]: ", len(sr.Channels)-1), len(sr.Channels)-1)
	}

	channel := sr.Channels[chanSelection]

	streamData, err := twitchClient().GetStreamData(channel.Id)
	if err != nil {
		return err
	}

	if streamData.Stream.Id == 0 {
		return fmt.Errorf("No online stream found for channel %s", ctx.Args()[0])
	}

	uris, err := twitchClient().GetStreamUrls(streamData.Stream.Channel.Name)
	if err != nil {
		return err
	}
	fmt.Printf("\n%s playing %s for %d viewers: %s\n\n",
		streamData.Stream.Channel.DisplayName,
		streamData.Stream.Game,
		streamData.Stream.Viewers,
		streamData.Stream.Channel.Status,
	)
	for i, uri := range uris {
		fmt.Printf("[%d]: %s (%s, %dkbps)\n", i, uri.Resolution, uri.Quality, uri.Bandwidth/1024)
	}

	streamSelection := getNumericInput(fmt.Sprintf("\nSelect stream format: [0-%d]: ", len(uris)-1), len(uris)-1)

	fmt.Printf("Loading %s %s (%s)...\n", channel.Name, uris[streamSelection].Resolution, uris[streamSelection].Quality)
	if err := mediaPlayer().LoadFromUrl(uris[streamSelection].URI); err != nil {
		return err
	}
	fmt.Printf("Playing %s %s (%s)...\n", channel.Name, uris[streamSelection].Resolution, uris[streamSelection].Quality)
	if err := mediaPlayer().Play(); err != nil {
		return err
	}
	if ctx.Bool("fullscreen") {
		fmt.Println("Entering fullscreen...")
		if err := mediaPlayer().EnterFullscreen(); err != nil {
			return err
		}
	}

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigchan
	switch sig {
	case syscall.SIGABRT:
		fmt.Println("Stream aborted!")
	case syscall.SIGINT:
		fmt.Println("Stream interrupted!")
	case syscall.SIGTERM:
		fmt.Println("Stream terminated!")
	}

	return nil
}
