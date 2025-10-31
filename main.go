package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"aleesa-telegram-go/internal/log"

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

	configJSONPath := fmt.Sprintf("%s/data/config.json", filepath.Dir(executablePath))

	locations := []string{
		"~/.aleesa-telegram-go.json",
		"~/aleesa-telegram-go.json",
		"/etc/aleesa-telegram-go.json",
		configJSONPath,
	}

	for _, location := range locations {
		config, err = parseConfig(location)

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
	if config.Log != "" {
		logfile, err = os.OpenFile(config.Log, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

		if err != nil {
			log.Errorf("Unable to open log file %s: %s", config.Log, err)
			os.Exit(1)
		}
	}

	log.Init(config.Loglevel, logfile)

	// Main context
	var ctx = context.Background()

	// Инициализируем redis-клиента.
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Redis.Server, config.Redis.Port),
	})

	log.Debugf("Lazy connect() to redis at %s:%d", config.Redis.Server, config.Redis.Port)
	subscriber = redisClient.Subscribe(ctx, config.Redis.MyChannel)
	redisMsgChan := subscriber.Channel()

	// Самое время поставить траппер сигналов.
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go sigHandler()
	go telega(config)

	// Обработчик событий от редиски.
	for msg := range redisMsgChan {
		if shutdown {
			continue
		}

		redisMsgParser(msg.Payload)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
