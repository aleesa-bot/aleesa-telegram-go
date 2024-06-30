package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/hjson/hjson-go"
	log "github.com/sirupsen/logrus"
)

// readConfig считывает и парсит конфиг-файл.
func readConfig() {
	configLoaded := false
	executablePath, err := os.Executable()

	if err != nil {
		log.Errorf("Unable to get current executable path: %s", err)
	}

	configJSONPath := fmt.Sprintf("%s/data/config.json", filepath.Dir(executablePath))

	locations := []string{
		"~/.aleesa-telegram-go.json",
		"~/aleesa-telegram-go.json",
		"/etc/aleesa-telegram-go.json",
		configJSONPath,
	}

	for _, location := range locations {
		fileInfo, err := os.Stat(location)

		// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём.
		if err != nil {
			continue
		}

		// Конфиг-файл длинноват для конфига, попробуем следующего кандидата.
		if fileInfo.Size() > 65535 {
			log.Warnf("Config file %s is too long for config, skipping", location)

			continue
		}

		buf, err := os.ReadFile(location)

		// Не удалось прочитать, попробуем следующего кандидата.
		if err != nil {
			log.Warnf("Skip reading config file %s: %s", location, err)

			continue
		}

		// Исходя из документации, hjson какбы умеет парсить "кривой" json, но парсит его в map-ку.
		// Интереснее на выходе получить структурку: то есть мы вначале конфиг преобразуем в map-ку, затем эту map-ку
		// сериализуем в json, а потом json преврщааем в стркутурку. Не очень эффективно, но он и не часто требуется.
		var (
			sampleConfig myConfig
			tmp          map[string]interface{}
		)

		err = hjson.Unmarshal(buf, &tmp)

		// Не удалось распарсить - попробуем следующего кандидата.
		if err != nil {
			log.Warnf("Skip parsing config file %s: %s", location, err)

			continue
		}

		tmpjson, err := json.Marshal(tmp)

		// Не удалось преобразовать map-ку в json.
		if err != nil {
			log.Warnf("Skip parsing config file %s: %s", location, err)

			continue
		}

		if err := json.Unmarshal(tmpjson, &sampleConfig); err != nil {
			log.Warnf("Skip parsing config file %s: %s", location, err)

			continue
		}

		// Валидируем значения из конфига.
		if sampleConfig.Redis.Server == "" {
			sampleConfig.Redis.Server = "localhost"

			log.Infof("Redis server is not defined in config, using localhost")
		}

		if sampleConfig.Redis.Port == 0 {
			sampleConfig.Redis.Port = 6379
		}

		if sampleConfig.Redis.Channel == "" {
			log.Errorf("Channel field in config file %s must be set", location)
			os.Exit(1)
		}

		if sampleConfig.Redis.MyChannel == "" {
			log.Errorf("My_channel field in config file %s must be set", location)
			os.Exit(1)
		}

		if sampleConfig.Loglevel == "" {
			sampleConfig.Loglevel = "info"
		}

		// sampleConfig.Log = "" if not set

		if sampleConfig.Csign == "" {
			log.Errorf("Csign field in config file %s must be set", location)
			os.Exit(1)
		}

		if sampleConfig.ForwardsMax == 0 {
			sampleConfig.ForwardsMax = forwardMax
		}

		if sampleConfig.DataDir == "" {
			log.Errorf("Data_dir field in config file %s must be set", location)
			os.Exit(1)
		}

		config = sampleConfig
		configLoaded = true

		log.Infof("Using %s as config file", location)

		break
	}

	if !configLoaded {
		log.Error("Config was not loaded! Refusing to start.")
		os.Exit(1)
	}
}

// sigHandler хэндлер сигналов закрывает все бд, все сетевые соединения и сваливает из приложения.
func sigHandler() {
	var err error

	for {
		var s = <-sigChan
		switch s {
		case syscall.SIGINT:
			log.Infoln("Got SIGINT, quitting")
		case syscall.SIGTERM:
			log.Infoln("Got SIGTERM, quitting")
		case syscall.SIGQUIT:
			log.Infoln("Got SIGQUIT, quitting")

		// Заходим на новую итерацию, если у нас "неинтересный" сигнал.
		default:
			continue
		}

		// Чтобы не срать в логи ошибками от редиски, проставим shutdown state приложения в true.
		shutdown = true

		// Отпишемся от всех каналов и закроем коннект к редиске.
		if err = subscriber.Unsubscribe(ctx); err != nil {
			log.Errorf("Unable to unsubscribe from redis channels cleanly: %s", err)
		} else {
			log.Debug("Unsubscribe from all redis channels")
		}

		if err = subscriber.Close(); err != nil {
			log.Errorf("Unable to close redis connection cleanly: %s", err)
		} else {
			log.Debug("Close redis connection")
		}

		// close all telego stuff.
		resp, err := tg.LogOut()

		if err != nil {
			log.Errorf("Unable to send LogOut() to Telegram Bot API server")
		}

		if !resp.Ok {
			log.Errorf("Telegram Bot API returns an error on LogOut(): %d, %s", resp.ErrorCode, resp.Description)
		}

		resp, err = tg.Close()

		if err != nil {
			log.Errorf("Unable to send Close() to Telegram Bot API server")
		}

		if !resp.Ok {
			log.Errorf("Telegram Bot API returns an error on Close(): %d, %s", resp.ErrorCode, resp.Description)
		}

		if len(settingsDB) > 0 {
			log.Debug("Closing runtime bot settings db")

			for _, db := range settingsDB {
				_ = db.Close()
			}
		}

		os.Exit(0)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
