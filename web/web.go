package web

import (
	"ShortLink/config"
)

func Start(cfg *config.Config) error {
	var err error
	go func() {
		err = startHTTPServer(cfg)
	}()
	return err
}
