package main

import (
	"encoding/json"
	"strconv"

	"github.com/NicoNex/echotron/v3"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// telego msg parser

// redisMsgParser парсит json-чики прилетевшие из REDIS-ки, причём, json-чики должны быть относительно валидными.
func redisMsgParser(msg string) {
	if shutdown {
		// Если мы завершаем работу программы, то нам ничего обрабатывать не надо.
		return
	}

	var j rMsg

	log.Debugf("Incoming raw json: %s", msg)

	if err := json.Unmarshal([]byte(msg), &j); err != nil {
		log.Warnf("Unable to to parse message from redis channel: %s", err)

		return
	}

	// Validate our j
	if exist := j.From; exist == "" {
		log.Warnf("Incorrect msg from redis, no from field: %s", msg)

		return
	}

	if exist := j.Chatid; exist == "" {
		log.Warnf("Incorrect msg from redis, no chatid field: %s", msg)

		return
	}

	if exist := j.Userid; exist == "" {
		log.Warnf("Incorrect msg from redis, no userid field: %s", msg)

		return
	}

	if exist := j.Message; exist == "" {
		log.Warnf("Incorrect msg from redis, no message field: %s", msg)

		return
	}

	if exist := j.Plugin; exist == "" {
		log.Warnf("Incorrect msg from redis, no plugin field: %s", msg)

		return
	}

	if exist := j.Mode; exist == "" {
		log.Warnf("Incorrect msg from redis, no mode field: %s", msg)

		return
	}

	// j.Misc.Answer может и не быть, тогда ответа на такое сообщение не будет
	if j.Misc.Answer == 0 {
		log.Debug("Field Misc->Answer = 0, skipping message")

		return
	}

	// j.Misc.BotNick тоже можно не передавать, тогда будет записана пустая строка
	// j.Misc.CSign если нам его не передали, возьмём значение из конфига. (но по идее нам оно тут не нужнО.)
	if exist := j.Misc.Csign; exist == "" {
		j.Misc.Csign = config.Csign
	}

	// j.Misc.FwdCnt если нам его не передали, то будет 0
	if exist := j.Misc.Fwdcnt; exist == 0 {
		j.Misc.Fwdcnt = 1
	}

	// j.Misc.GoodMorning может быть 1 или 0, по-умолчанию 0
	// j.Misc.MsgFormat может быть 1 или 0, по-умолчанию 0
	// j.Misc.Username можно не передавать, тогда будет пустая строка

	// Отвалидировались, теперь вернёмся к нашим баранам.

	var opts *echotron.MessageOptions

	chatid, err := strconv.ParseInt(j.Chatid, 10, 64)

	if err != nil {
		log.Errorf("unable to parse message from redis, incorrect chatid field: %s", err)

		return
	}

	resp, err := tg.SendMessage(j.Message, chatid, opts)

	if err != nil || !resp.Ok {
		// TODO: поддержать миграцию группы в супергруппу, поддержать вариант, когда бот замьючен.
		// Красиво оформить ошибку, с полями итд, как tracedump, только ошибка.
		// N.B. тут может быть сообщение о том, что группа превратилась в супергруппу, или что бот не имеет прав писать
		// сообщения в чятик. Это надо бы отслеживать и хэндлить.
		log.Errorf("Unable to send message to telegram api: %s", err)
		log.Errorf("Response dump: %s", spew.Sdump(resp))
	}
}

// telegramMsgParser парсит ивент, прилетевший из bot api.
func telegramMsgParser(msg *echotron.Update) {
	if shutdown {
		// Если мы завершаем работу программы, то нам ничего обрабатывать не надо.
		return
	}

	// Сообщение о том, что этот чятик изменился, например превратился в супергруппу.
	if msg.Message.MigrateFromChatID < 0 && msg.Message.MigrateToChatID < 0 {
		// TODO: поддержать миграцию настроек чата га новый chatId.
		// TODO: подумать, что можно сделать с настройками бэкэндов. Кажись, ничего, но надо глянуть.
	}

	// Люди пришли в чят.
	if msg.Message.NewChatMembers != nil {
		// TODO: проверить, надо ли приветсвовать, выбрать одну из приветственных фраз, потянуть время, типа, мы пишем
		// текст и поприветсвовать вновьприбывшего
	}

	// Человек ушёл из чата.
	if msg.Message.LeftChatMember != nil {
		// TODO: проверить, надо ли прощаться с ушедшим участником чата и попрощаться, если надо
	}

	var rmsg rMsg

	switch msg.Message.Chat.Type {
	case "private":
		// handle private messages
		rmsg.Misc.Answer = 1
		// если фраза является командой - засылаем её в парсер команд

		// Засылаем фразу в misc-канал (в роутер)
	case "group", "supergroup":
		// handle public (super)group messages
		// Здесь работает цензор, если он включён.
		// Здесь же, если человек ответил на авто-приветствие, это надо обнаруживать и пропускать.
		// Здесь же, если фраза менее 3-х букв, просто игнорируем её
		// Здесь же, если фраза является командой - засылаем её в парсер команд
		// Здесь же, если фраза обращена к боту (ник или имя), выставляем флажок, что надо ответить
		// Здесь же, пытаемся убрать из фразы ник или имя бота
		// Засылаем фразу в misc-канал (в роутер)
	case "channel":
		// Тут мы ничего сделать не можем.
	} //nolint:wsl
}

func cmdParser(cmd string) error {
	var err error

	return err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
