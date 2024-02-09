package main

import (
	"github.com/NicoNex/echotron/v3"
)

// telega основная горутинка, реализующая бота.
func telega(c myConfig) {
	tg = echotron.NewAPI(c.Telegram.Token)

	for u := range echotron.PollingUpdates(c.Telegram.Token) {
		telegramMsgParser(u)
	}
}
