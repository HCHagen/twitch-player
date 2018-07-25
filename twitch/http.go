package twitch

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

type StatusError struct {
	error
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func newDialer() func(context.Context, string, string) (net.Conn, error) {
	return (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext
}

func newTransport(tlsConfig *tls.Config, dialer func(context.Context, string, string) (net.Conn, error)) *http.Transport {
	return &http.Transport{
		DialContext:           dialer,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		TLSClientConfig: tlsConfig,
	}
}

func newHttpClient(timeout time.Duration, tlsConfig *tls.Config) *http.Client {
	return &http.Client{
		Transport: newTransport(tlsConfig, newDialer()),
		Timeout:   timeout,
	}
}

func newHttp2Client(timeout time.Duration, tlsConfig *tls.Config) (*http.Client, error) {
	c := &http.Client{
		Transport: newTransport(tlsConfig, newDialer()),
		Timeout:   timeout,
	}
	return c, http2.ConfigureTransport(c.Transport.(*http.Transport))
}
