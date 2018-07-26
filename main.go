package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/urfave/cli"

	"github.com/HCHagen/twitch-player/player"
	"github.com/HCHagen/twitch-player/twitch"
)

var (
	DefaultTwitchHttpTimeout = 20 * time.Second

	mediaPlayerCloser func() error = func() error {
		return nil
	}

	aout, vout  string
	mediaPlayer func() player.Player = func() func() player.Player {
		var p player.Player
		var err error
		return func() player.Player {
			if p != nil {
				return p
			}

			p, err = player.NewVlcPlayer(aout, vout)
			if err != nil {
				fmt.Printf("Error initializing media player: %s\n", err.Error())
				os.Exit(1)
			}
			mediaPlayerCloser = p.Close

			return p
		}
	}()

	twitchClient func() twitch.Client = func() func() twitch.Client {
		var c twitch.Client
		var err error
		return func() twitch.Client {
			if c != nil {
				return c
			}

			c, err = twitch.NewTwitchClient(appClientId, DefaultTwitchHttpTimeout)
			if err != nil {
				fmt.Printf("Error initializing twitch client: %s\n", err.Error())
				os.Exit(1)
			}

			return c
		}
	}()
)

func main() {
	app := cli.NewApp()
	app.Name = "twitch-player"
	app.Author = fmt.Sprintf("Built by: %s at %s", appBuildUser, buildTime())
	app.Version = appVersion
	app.Usage = appUsage
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr

	app.Before = func(ctx *cli.Context) error {
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name: "stream",
			Before: func(ctx *cli.Context) error {
				aout = ctx.String("aout")
				vout = ctx.String("vout")
				return nil
			},
			Usage:  "Play stream from channel",
			Action: onStream,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "fullscreen,f",
					Usage: "Run in fullscreen",
				},
				cli.StringFlag{
					Name:  "aout,a",
					Usage: "Audio output device",
				},
				cli.StringFlag{
					Name:  "vout,v",
					Usage: "Video output device",
				},
			},
		},
		{
			Name:   "list",
			Usage:  "List stream channels",
			Action: onList,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "game,g",
					Usage: "Game to list channels for",
				},
				cli.IntFlag{
					Name:  "number,n",
					Usage: "Number of channels to list",
				},
			},
		},
		{
			Name:   "search",
			Usage:  "Search for channels",
			Action: onSearch,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "number,n",
					Usage: "Number of results to list",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	if err := mediaPlayerCloser(); err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
}

func onStream(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("Please provide a channel name")
	}

	streamData, err := twitchClient().GetStreamData(ctx.Args()[0])
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
		streamData.Stream.Channel.Game,
		streamData.Stream.Viewers,
		streamData.Stream.Channel.Status,
	)
	for i, uri := range uris {
		fmt.Printf("[%d]: %s (%s, %dkbps)\n", i, uri.Resolution, uri.Quality, uri.Bandwidth/1024)
	}

	reader := bufio.NewReader(os.Stdin)
	selection := -1
	for selection < 0 {
		fmt.Printf("\nSelect stream format: [0-%d]: ", len(uris)-1)
		if inp, _ := reader.ReadString('\n'); inp != "\n" {
			if i, err := strconv.Atoi(strings.TrimSpace(inp)); err != nil {
				fmt.Printf("%s is not a valid number...\n", inp)
				continue
			} else {
				if i < 0 || i > len(uris)-1 {
					fmt.Printf("%s is not in range [0-%d]...\n", inp, len(uris)-1)
					continue
				}
				selection = i
			}
		} else {
			selection = 0
		}
	}

	fmt.Printf("Loading %s %s (%s)...\n", ctx.Args()[0], uris[selection].Resolution, uris[selection].Quality)
	if err := mediaPlayer().LoadFromUrl(uris[selection].URI); err != nil {
		return err
	}
	fmt.Printf("Playing %s %s (%s)...\n", ctx.Args()[0], uris[selection].Resolution, uris[selection].Quality)
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

func onList(ctx *cli.Context) error {
	return fmt.Errorf("Not implemented yet...")
}

func onSearch(ctx *cli.Context) error {
	return fmt.Errorf("Not implemented yet...")
}
