package main

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

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
