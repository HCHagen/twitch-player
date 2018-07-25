package player

type Player interface {
	LoadFromUrl(url string) error
	LoadFromFile(path string) error

	Play() error
	Stop() error

	EnterFullscreen() error
	ExitFullscreen() error

	Reset() error
	Close() error
}
