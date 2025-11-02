package tg

import (
	"context"
	"os"

	"github.com/NicoNex/echotron/v3"
	"github.com/carlescere/scheduler"
	"github.com/cockroachdb/pebble"
	"github.com/go-redis/redis/v8"
)

var (
	// Config - это у нас глобальная штука.
	Config MyConfig

	// forwardMax is used to break circular message forwarding we must set some sane default, it can be overridden via config.
	forwardMax int64 = 5

	// RedisClient объектик клиента-редиски.
	RedisClient *redis.Client

	// Subscriber is redis PubSub object.
	Subscriber *redis.PubSub

	// Ctx is main context.
	Ctx = context.Background()

	// Shutdown ставится в true, если мы получили сигнал на выключение.
	Shutdown = false

	// SigChan канал, в который приходят уведомления для хэндлера сигналов от траппера сигналов.
	SigChan = make(chan os.Signal, 1)

	// settingsDB мапка с открытыми дескрипторами баз с настройками.
	settingsDB = make(map[string]*pebble.DB)

	// chatListDB объектик *pebble.DB с базой, в которой лежит список чатов.
	chatListDB *pebble.DB

	// chatListDBName имя базы данных со списком чатов.
	chatListDBName = "chat_list_db"

	// chatList слайс со списоком чатов.
	chatList = []string{}

	// tg is var for telegram api.
	tg echotron.API

	// introduceGreet приветственные фразы для новых участников чата.
	introduceGreet = [...]string{
		"Дратути",
		"Дарована",
		"Доброе утро, день или вечер",
		"Добро пожаловать в наше скромное коммунити",
		"Наше вам с кисточкой тут, на канальчике",
	}

	// PeriodicJobs contains pointer to slice with sheduled job objects.
	PeriodicJobs = make([]*scheduler.Job, 0)
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
