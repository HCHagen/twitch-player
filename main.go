package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/hchagen/twitch-player/player"
	"github.com/hchagen/twitch-player/twitch"
)

var (
	DefaultTwitchHttpTimeout = 20 * time.Second
	DefaultChannelResultLen  = 5
	DefaultGameResultLen     = 25
	DefaultStreamResultLen   = 30

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
			Name:   "games",
			Usage:  "Display games",
			Action: onListGames,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "number,n",
					Usage: "Number of channels to list",
					Value: DefaultGameResultLen,
				},
			},
		},
		{
			Name:   "list",
			Usage:  "List stream channels",
			Action: onListStreams,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "featured,f",
					Usage: "List featured channels",
				},
				cli.StringFlag{
					Name:  "game,g",
					Usage: "Game to list channels for",
				},
				cli.IntFlag{
					Name:  "number,n",
					Usage: "Number of channels to list",
					Value: DefaultStreamResultLen,
				},
			},
		},
		{
			Name:   "search",
			Usage:  "Search for streams",
			Action: onSearch,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "channel,c",
					Usage: "Search for channels rather than streams",
				},
				cli.BoolFlag{
					Name:  "game,g",
					Usage: "Search for games rather than streams",
				},
				cli.IntFlag{
					Name:  "number,n",
					Usage: "Number of results to list",
					Value: DefaultChannelResultLen,
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

func printChannels(channels []twitch.Channel) {
	for _, channel := range channels {
		fmt.Printf("[%s] %s (id %d) last played: %s (%s)\n", channel.Name, channel.DisplayName, channel.Id, channel.Game, channel.Updated)
	}
}

func printGamesInfo(games []twitch.GameInfo) {
	for _, game := range games {
		fmt.Printf("[%s] %s has %d viewers on %d channels\n", game.Game.Name, game.Game.LocalizedName, game.Viewers, game.Channels)
	}
}

func printGames(games []twitch.Game) {
	for _, game := range games {
		fmt.Printf("[%s] %s (id %d)\n", game.Name, game.LocalizedName, game.Id)
	}
}

func printFeatured(featured []twitch.Featured) {
	for _, stream := range featured {
		fmt.Printf("[%s] %s playing %s for %d viewers: %s\n",
			stream.Stream.Channel.Name,
			stream.Stream.Channel.DisplayName,
			stream.Stream.Game,
			stream.Stream.Viewers,
			stream.Title,
		)
	}
}

func printStreams(streams []twitch.Stream) {
	for _, stream := range streams {
		fmt.Printf("[%s] %s playing %s for %d viewers: %s\n",
			stream.Channel.Name,
			stream.Channel.DisplayName,
			stream.Game,
			stream.Viewers,
			stream.Channel.Status,
		)
	}
}

func getNumericInput(prompt string, max int) int {
	reader := bufio.NewReader(os.Stdin)
	selection := -1
	for selection < 0 {
		fmt.Printf(prompt)
		if inp, _ := reader.ReadString('\n'); inp != "\n" {
			if i, err := strconv.Atoi(strings.TrimSpace(inp)); err != nil {
				fmt.Printf("%s is not a valid number...\n", inp)
				continue
			} else {
				if i < 0 || i > max {
					fmt.Printf("%s is not in range [0-%d]...\n", inp, max)
					continue
				}
				selection = i
			}
		} else {
			selection = 0
		}
	}

	return selection
}
