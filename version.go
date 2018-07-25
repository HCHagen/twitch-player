package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var appBuildTime = "0"
var appBuildUser = "unknown"
var appClientId = ""
var appVersion = "0.1.0+git"

var appUsage = `A twitch cli for browsing/playing streams`

func buildTime() time.Time {
	bt, err := strconv.Atoi(appBuildTime)
	if err != nil {
		fmt.Printf("Unable to parse app build time: %s\n", err.Error())
		os.Exit(1)
	}

	return time.Unix(int64(bt), 0)
}
