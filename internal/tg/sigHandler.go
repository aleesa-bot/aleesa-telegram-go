package tg

import (
	"os"
	"syscall"

	"aleesa-telegram-go/internal/log"
)

// SigHandler хэндлер сигналов закрывает все бд, все сетевые соединения и сваливает из приложения.
func SigHandler() {
	var err error

	for {
		var s = <-SigChan
		switch s {
		case syscall.SIGINT:
			log.Info("Got SIGINT, quitting")
		case syscall.SIGTERM:
			log.Info("Got SIGTERM, quitting")
		case syscall.SIGQUIT:
			log.Info("Got SIGQUIT, quitting")

		// Заходим на новую итерацию, если у нас "неинтересный" сигнал.
		default:
			continue
		}

		// Чтобы не срать в логи ошибками от редиски, проставим shutdown state приложения в true.
		Shutdown = true

		// Отпишемся от всех каналов и закроем коннект к редиске.
		if err = Subscriber.Unsubscribe(Ctx); err != nil {
			log.Errorf("Unable to unsubscribe from redis channels cleanly: %s", err)
		} else {
			log.Debug("Unsubscribe from all redis channels")
		}

		if err = Subscriber.Close(); err != nil {
			log.Errorf("Unable to close redis connection cleanly: %s", err)
		} else {
			log.Debug("Close redis connection")
		}

		// Quit all scheduled periodic jobs.
		for _, job := range PeriodicJobs {
			job.Quit <- true
		}

		resp, err := tg.Close()

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

		log.Debug("Closing known chat list db")
		_ = chatListDB.Close()

		os.Exit(0)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
