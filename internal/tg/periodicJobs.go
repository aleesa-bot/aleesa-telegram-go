package tg

import (
	"aleesa-telegram-go/internal/log"
	"encoding/json"
	"strconv"
)

var (
	// TidySettingsDB флашит на диск базу с настройками чатов. Рекомендуется делать это время от времени для поддержания
	// производительности базы на должном уровне.
	// Предполагается, что это происходит где-то раз в час, чаще смысла нет.
	TidySettingsDB = func() {
		if len(settingsDB) > 0 {
			for name, db := range settingsDB {
				log.Debugf("Flushing %s settings db", name)

				if err := db.Flush(); err != nil {
					log.Errorf("Unable to Flush() %s db: %s", name, err)
				}
			}
		}

		log.Debugf("Flushing %s", chatListDBName)

		if err := chatListDB.Flush(); err != nil {
			log.Errorf("Unable to Flush() %s: %s", chatListDBName, err)
		}
	}

	// SendGoodMorning отсылает сообщение "с добрым утром" (сейчас - фортунку). Предполагается, что оно прогоняется по
	// всем желающим чатам раз в сутки, с утра. Утро отсчитывается от настроек локали в системе.
	SendGoodMorning = func() {
		me, err := tg.GetMe()

		if err != nil {
			log.Errorf("Unable to send good mornings message, cannot get info about myself: %s", err)

			return
		}

		for _, chatID := range chatGroupList {
			if GetSetting(chatID, "FortuneMsg") == "1" {
				// Засылаем фразу в misc-канал (в роутер).
				var rmsg rMsg

				rmsg.Chatid = chatID
				rmsg.Userid = strconv.FormatInt(me.Result.ID, 10)
				rmsg.Message = "!f"
				rmsg.Mode = "public"
				rmsg.Plugin = "telegram"
				rmsg.From = "telegram"
				rmsg.Misc.Csign = Config.Csign
				rmsg.Misc.Fwdcnt = 1

				rmsg.Misc.Botnick = ConstructPartialUserUsername(me.Result)
				rmsg.Misc.Username = ConstructPartialUserUsername(me.Result)

				rmsg.Misc.GoodMorning = 1
				// Форматирование нужно только для вывода некоторых ответов на команды, команды мы ловим выше по тексту, так что
				// смело ставим тут 0.
				rmsg.Misc.Msgformat = 0

				data, err := json.Marshal(rmsg)

				if err != nil {
					log.Warnf("Unable to to serialize message for redis: %s", err)

					return
				}

				// Заталкиваем наш json в редиску.
				if err := RedisClient.Publish(Ctx, Config.Redis.Channel, data).Err(); err != nil {
					log.Warnf("Unable to send data to redis channel %s: %s", Config.Redis.Channel, err)
				} else {
					log.Debugf("Sent msg to redis channel %s: %s", Config.Redis.Channel, string(data))
				}
			}
		}
	}
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
