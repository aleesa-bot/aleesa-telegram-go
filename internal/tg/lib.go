package tg

import (
	"aleesa-telegram-go/internal/log"
	"os"
	"path/filepath"
)

// ReadConfig читает конфиг. Есть несколько преодопределённых локаций, откуда он пытается его вычитать.
// ./data/config.json, ~/.aleesa-telegram-go.json, ~/aleesa-telegram-go.json, ~/aleesa-telegram-go.json.
func ReadConfig() {
	var (
		err error
	)

	executablePath, err := os.Executable()

	if err != nil {
		log.Fatalf("Unable to get current executable path: %s", err)
	}

	configJSONPath := filepath.Dir(executablePath) + "/data/config.json"

	locations := []string{
		"~/.aleesa-telegram-go.json",
		"~/aleesa-telegram-go.json",
		"/etc/aleesa-telegram-go.json",
		configJSONPath,
	}

	for _, location := range locations {
		Config, err = parseConfig(location)

		if err == nil {
			break
		}

		log.Errorf("Unable to parse config at %s: %s", location, err)
	}

	if err != nil {
		os.Exit(1)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
