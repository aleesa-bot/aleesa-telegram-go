package tg

import (
	"aleesa-telegram-go/internal/log"
)

var (
	// TidySettingsDB flushes regularly pebble dbs with settings. Remmended action to keeb that db in a good shape.
	// Meant to be run about hourly or something like that.
	TidySettingsDB = func() {
		if len(settingsDB) > 0 {
			for name, db := range settingsDB {
				log.Debugf("Flushing %s settings db", name)

				if err := db.Flush(); err != nil {
					log.Errorf("Unable to Flush() %s db: %s", name, err)
				}
			}
		}
	}
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
