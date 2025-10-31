package tg

import (
	"context"
	"os"

	"github.com/NicoNex/echotron/v3"
	"github.com/cockroachdb/pebble"
	"github.com/go-redis/redis/v8"
)

// Config - это у нас глобальная штука.
var Config myConfig

// To break circular message forwarding we must set some sane default, it can be overridden via config.
var forwardMax int64 = 5

// Объектики клиента-редиски.
var RedisClient *redis.Client
var Subscriber *redis.PubSub

// Main context.
var ctx = context.Background()

// Ставится в true, если мы получили сигнал на выключение.
var Shutdown = false

// Канал, в который приходят уведомления для хэндлера сигналов от траппера сигналов.
var SigChan = make(chan os.Signal, 1)

// Мапка с открытыми дескрипторами баз с настройками.
var settingsDB = make(map[string]*pebble.DB)

var tg echotron.API

// Приветсвенные фразы для новых участников чата.
var introduceGreet = [...]string{
	"Дратути",
	"Дарована",
	"Доброе утро, день или вечер",
	"Добро пожаловать в наше скромное коммунити",
	"Наше вам с кисточкой тут, на канальчике",
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
