package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func onListGames(ctx *cli.Context) error {
	games, err := twitchClient().GetGameList(ctx.Int("number"))
	if err != nil {
		return err
	}

	fmt.Printf("Listing top %d (out of %d) games:\n\n", len(games.Games), games.Total)
	printGamesInfo(games.Games)
	fmt.Println("")

	return nil
}

func listFeatured(num int) error {
	featured, err := twitchClient().GetFeaturedList(num)
	if err != nil {
		return err
	}

	fmt.Printf("Listing top %d featured streams:\n\n", len(featured.Featured))
	printFeatured(featured.Featured)

	return nil
}

func listStreams(game string, num int) error {
	streams, err := twitchClient().GetStreamList(game, num)
	if err != nil {
		return err
	}

	fmt.Printf("Listing top %d streamers (out of %d):\n\n", len(streams.Streams), streams.Total)
	printStreams(streams.Streams)

	return nil
}

func onListStreams(ctx *cli.Context) error {
	var err error
	if ctx.Bool("featured") {
		err = listFeatured(ctx.Int("number"))
	} else {
		err = listStreams(ctx.String("game"), ctx.Int("number"))
	}
	if err != nil {
		return err
	}
	fmt.Println("")

	return nil
}
