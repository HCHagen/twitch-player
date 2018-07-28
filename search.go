package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func channelSearch(channel string, num int) error {
	searchResult, err := twitchClient().GetChannelSearch(channel, num)
	if err != nil {
		return err
	}

	fmt.Printf("%d results - displaying first %d:\n\n", searchResult.Total, len(searchResult.Channels))
	printChannels(searchResult.Channels)

	return nil
}

func gameSearch(game string) error {
	searchResult, err := twitchClient().GetGameSearch(game)
	if err != nil {
		return err
	}

	fmt.Printf("%d results - displaying first %d:\n\n", len(searchResult.Games), len(searchResult.Games))
	printGames(searchResult.Games)

	return nil
}

func streamSearch(channel string, num int) error {
	searchResult, err := twitchClient().GetStreamSearch(channel, num)
	if err != nil {
		return err
	}

	fmt.Printf("%d results - displaying first %d:\n\n", searchResult.Total, len(searchResult.Streams))
	printStreams(searchResult.Streams)

	return nil
}

func onSearch(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("Please provide a search query")
	}

	var err error
	if ctx.Bool("channel") {
		err = channelSearch(ctx.Args()[0], ctx.Int("number"))
	} else if ctx.Bool("game") {
		err = gameSearch(ctx.Args()[0])
	} else {
		err = streamSearch(ctx.Args()[0], ctx.Int("number"))
	}
	if err != nil {
		return err
	}
	fmt.Println("")

	return nil
}
