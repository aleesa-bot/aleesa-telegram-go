// Package main is the main package of aleesa-telegram-go.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"aleesa-telegram-go/internal/log"
	"aleesa-telegram-go/internal/tg"

	"github.com/carlescere/scheduler"
	"github.com/go-redis/redis/v8"
)

// init производит некоторую инициализацию перед запуском main().
func init() {
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
		tg.Config, err = tg.ParseConfig(location)

		if err == nil {
			break
		}

		log.Errorf("Unable to parse config at %s: %s", location, err)
	}

	if err != nil {
		os.Exit(1)
	}
}

func main() {
	var (
		logfile *os.File
		err     error
	)

	// Откроем лог и скормим его логгеру.
	if tg.Config.Log != "" {
		logfile, err = os.OpenFile(tg.Config.Log, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

		if err != nil {
			log.Errorf("Unable to open log file %s: %s", tg.Config.Log, err)
			os.Exit(1)
		}
	}

	log.Init(tg.Config.Loglevel, logfile)

	// Инициализируем redis-клиента.
	tg.RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", tg.Config.Redis.Server, tg.Config.Redis.Port),
	})

	log.Debugf("Lazy connect() to redis at %s:%d", tg.Config.Redis.Server, tg.Config.Redis.Port)

	tg.Subscriber = tg.RedisClient.Subscribe(tg.Ctx, tg.Config.Redis.MyChannel)

	redisMsgChan := tg.Subscriber.Channel()

	// Самое время поставить траппер сигналов.
	signal.Notify(tg.SigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go tg.SigHandler()
	go tg.Telega(tg.Config)

	// Периодически обслуживаем базы с настройками чатов.
	if job, err := scheduler.Every(1).Hours().NotImmediately().Run(tg.TidySettingsDB); err != nil {
		log.Errorf("Unable to schedule periodic settings db flush: %s", err)
	} else {
		tg.PeriodicJobs = append(tg.PeriodicJobs, job)
	}

	// Отправлем "доброе утро" каждое утро для всех чатов, которые на это подписались.
	if job, err := scheduler.Every().Day().At("8:10").Run(tg.SendGoodMorning); err != nil {
		log.Errorf("Unable to schedule send good morning task: %s", err)
	} else {
		tg.PeriodicJobs = append(tg.PeriodicJobs, job)
	}

	// Обработчик событий от редиски.
	for msg := range redisMsgChan {
		if tg.Shutdown {
			continue
		}

		tg.RedisMsgParser(msg.Payload)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
