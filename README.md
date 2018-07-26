# twitch-player
A cli twitch-player with (hopefully) embedded VLC. So far experimental work in progress.

# install
apt install libvlc-dev
go get github.com/HCHagen/twitch-player
cd ${GOPATH}/src/github.com/HCHagen/twitch-player
make

# usage
twitch-player stream "channelname"

Depends on VLC so far. May add support for streaming to fifo pipe and omxplayer playing from said pipe.
