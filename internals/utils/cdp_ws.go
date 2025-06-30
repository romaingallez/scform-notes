package utils

import (
	"context"

	"github.com/gobwas/ws"
	log "github.com/sirupsen/logrus"
)

func TestChromeDevWS(wsURL string) (err error) {

	log.Println("Testing Chrome Dev WS:", wsURL)

	conn, _, _, err := ws.Dial(context.Background(), wsURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil

}
