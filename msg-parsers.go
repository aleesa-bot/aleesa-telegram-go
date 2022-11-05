package main

import (
	"encoding/json"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// telego msg parser

// redisMsgParser парсит json-чики прилетевшие из REDIS-ки, причём, json-чики должны быть относительно валидными
func redisMsgParser(msg string) {
	if shutdown {
		// Если мы завершаем работу программы, то нам ничего обрабатывать не надо
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
	// j.Misc.CSign если нам его не передали, возьмём значение из конфига
	if exist := j.Misc.Csign; exist == "" {
		j.Misc.Csign = config.Csign
	}

	// j.Misc.FwdCnt если нам его не передали, то будет 0
	if exist := j.Misc.Fwdcnt; exist == 0 {
		j.Misc.Fwdcnt = 1
	}

	// j.Misc.GoodMorning может быть быть 1 или 0, по-умолчанию 0
	// j.Misc.MsgFormat может быть быть 1 или 0, по-умолчанию 0
	// j.Misc.Username можно не передавать, тогда будет пустая строка

	// Отвалидировались, теперь вернёмся к нашим баранам.
	lines := regexp.MustCompile("\r?\n").Split(j.Message, -1)

	for _, message := range lines {
		if userModeIsOped(j.Chatid, ircClient.GetNick()) || userModeIsVoiced(j.Chatid, ircClient.GetNick()) {
			imChanUnrestricted <- iMsg{ChatId: j.Chatid, Text: message}
		} else {
			imChan <- iMsg{ChatId: j.Chatid, Text: message}
		}
	}

	return
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
