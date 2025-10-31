package tg

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"

	"aleesa-telegram-go/internal/log"

	"github.com/NicoNex/echotron/v3"
	"github.com/davecgh/go-spew/spew"
)

// telegramMsgParser парсит ивент, прилетевший из bot api.
func telegramMsgParser(msg *echotron.Update) {
	if Shutdown {
		// Если мы завершаем работу программы, то нам ничего обрабатывать не надо.
		return
	}

	// Сообщение о том, что этот чятик изменился, например, превратился в супергруппу.
	if (msg.Message.MigrateFromChatID < 0) && (msg.Message.MigrateToChatID < 0) { //nolint: revive,staticcheck
		// TODO: поддержать миграцию настроек чата на новый chatId.
		// TODO: подумать, что можно сделать с настройками бэкэндов. Кажись, ничего, но надо глянуть.

		return
	}

	// Здесь начинается парсинг, подразумевающий какие-то сообщения бота в чат, поэтому с этой точки можно/нужно
	// получить информацию о самом боте.
	// TODO: пораскинуть, можно ли как-то кэшировать эту информацию, чтобы не дёргать api по чём зря.
	var (
		rmsg         rMsg
		errorOccured bool
	)

	me, err := tg.GetMe()

	if err != nil {
		log.Errorf("Unable to get info about myself, %s", err)

		errorOccured = true
	}

	// Типа, при некоторых обстоятельствах, мы можем получить более внятное сообщение. Но это надо проверять.
	if !me.Ok {
		log.Errorf("Unable to get info about myself: %d %s", me.ErrorCode, me.Description)

		errorOccured = true
	}

	if errorOccured {
		return
	}

	// Люди пришли в чят.
	if msg.Message.NewChatMembers != nil { //nolint: revive,staticcheck
		// TODO: проверить, надо ли приветствовать, выбрать одну из приветственных фраз, потянуть время, типа, мы пишем
		// текст и поприветствовать вновь прибывшего.
		if GetSetting(fmt.Sprintf("%d", msg.Message.Chat.ID), "GreetMsg") == "1" {
			// Пользователи, которых мы приветсвуем, одной строкой.
			var users string

			if len(msg.Message.NewChatMembers) > 1 {
				lastOne := msg.Message.NewChatMembers[len(msg.Message.NewChatMembers)-1]
				theRest := msg.Message.NewChatMembers[:len(msg.Message.NewChatMembers)-1]
				var usersSlice []string

				for _, user := range theRest {
					usersSlice = append(usersSlice, ConstructTelegramHighlightName(user))
				}

				users = strings.Join(usersSlice, ", ")
				users += fmt.Sprintf(" and %s", ConstructTelegramHighlightName(lastOne))
			} else {
				users = ConstructTelegramHighlightName(msg.Message.NewChatMembers[0])
			}

			phrase := fmt.Sprintf(
				"%s, %s. Представьтес, пожалуйста, и расскажите, что вас сюда привело.",
				introduceGreet[rand.IntN(len(introduceGreet))],
				users,
			)

			// Потянуть время, изобразить настоящего человека.
			resp, err := tg.SendChatAction(
				echotron.Typing,
				msg.Message.From.ID,
				&echotron.ChatActionOptions{MessageThreadID: int(msg.Message.ThreadID)},
			)

			// TODO: вероятно, ошибку надо засылать в обработчик ошибок, так как в ней может быть миграция чятика или
			// что-то, что тоже надо хэндлить.
			if err != nil || !resp.Ok {
				log.Errorf(
					"unable to send message to telegram api: %s\nResponse dump: %s",
					err,
					spew.Sdump(resp),
				)
			}

			resp1, err := tg.SendMessage(
				phrase,
				msg.Message.Chat.ID,
				&echotron.MessageOptions{ParseMode: "MarkdownV2", MessageThreadID: int64(msg.Message.ThreadID)},
			)

			// TODO: вероятно, ошибку надо засылать в обработчик ошибок, так как в ней может быть миграция чятика или
			// что-то, что тоже надо хэндлить.
			if err != nil || !resp1.Ok {
				log.Errorf(
					"unable to send message to telegram api: %s\nResponse dump: %s",
					err,
					spew.Sdump(resp1),
				)
			}
		}

		return
	}

	// Человек ушёл из чата.
	if msg.Message.LeftChatMember != nil {
		if GetSetting(fmt.Sprintf("%d", msg.Message.Chat.ID), "GoodbyeMsg") == "1" {
			user := ConstructUserFirstLastName(msg.Message.From)

			goodbye := fmt.Sprintf("Прощаемся с %s", user)

			opts := echotron.MessageOptions{DisableNotification: true, ParseMode: "MarkdownV2"}

			if msg.Message.ThreadID != 0 {
				opts.MessageThreadID = int64(msg.Message.ThreadID)
			}

			if res, err := tg.SendMessage(goodbye, msg.Message.Chat.ID, &opts); err != nil {
				log.Errorf("Unable to send message to telegram api: %s", err)
			} else if !res.Ok {
				log.Errorf("Unable to send message to telegram api: %s", res.Description)
			}
		}

		return
	}

	// Если сообщение было зацензурено, то дальше его обрабатывать не надо.
	if Censor(msg) {
		return
	}

	switch msg.Message.Chat.Type {
	case "private":
		// Всё что можно было похэндлить, что не содержало текста, считаем, что обработали.
		if msg.Message.Text == "" {
			return
		}

		// Считаем, что приватное сообщение всегда нуждается в ответе, иначе, зачем писать боту?
		rmsg.Misc.Answer = 1

		// Засылаем фразу в misc-канал (в роутер).
		rmsg.Chatid = fmt.Sprintf("%d", msg.Message.Chat.ID)
		rmsg.Userid = fmt.Sprintf("%d", msg.Message.From.ID)
		rmsg.Message = msg.Message.Text
		rmsg.Mode = "private"
		rmsg.Plugin = "telegram"
		rmsg.From = "telegram"
		rmsg.Misc.Csign = Config.Csign
		rmsg.Misc.Fwdcnt = 1
		rmsg.Misc.Botnick = ConstructPartialUserUsername(me.Result)
		rmsg.Misc.Username = ConstructPartialUserUsername(msg.Message.From)
		rmsg.Misc.GoodMorning = 0
		rmsg.Misc.Msgformat = 0

		// Форматировать что-то всё-таки надо, но благо, это только команды, поэтому попробуем распознать сообщение как
		// команду.
		r, err := cmdParser(me, msg)

		if err != nil {
			log.Infof(
				"Unable to parse message as command in private conversation with %s, message was: %s",
				ConstructFullUserName(msg.Message.From),
				msg.Message.Text,
			)

			return
		}

		// Если сообщение не было опознанно как команда, засылаем его "как есть" туда, в качель.
		if !r {
			data, err := json.Marshal(rmsg)

			if err != nil {
				log.Warnf("Unable to to serialize message for redis: %s", err)

				return
			}

			// Заталкиваем наш json в редиску.
			if err := RedisClient.Publish(ctx, Config.Redis.Channel, data).Err(); err != nil {
				log.Warnf("Unable to send data to redis channel %s: %s", Config.Redis.Channel, err)
			} else {
				log.Debugf("Sent msg to redis channel %s: %s", Config.Redis.Channel, string(data))
			}
		}

	case "group", "supergroup":
		// Обрабатываем сообщения для групп и супергрупп. По сути это примерно одинаковые вещи/сущности.
		// Цензор удаляет сообщения из чятика, если они "неправильные" - от имени других каналов, аудиосообщения,
		// содержат картинки, видео аудио итп, это настраивается через команду !admin censor.
		if GetSetting(fmt.Sprintf("%d", msg.Message.Chat.ID), "censor") == "1" {
			// Если цензор отработал, больше делать нечего, сообщения уже нету.
			if Censor(msg) {
				return
			}
		}

		// Если фраза менее 3-х букв, просто игнорируем её.
		if len(msg.Message.Text) <= 3 {
			return
		}

		// TODO: если человек ответил на авто-приветствие, это надо обнаруживать и пропускать.

		// Pасылаем сообщение в парсер команд.
		res, err := cmdParser(me, msg)

		if err != nil {
			log.Errorf(
				"Unable to parse message from telegram api as a command. Message was: %s, error: %s",
				msg.Message.Text,
				err,
			)
		}

		// Если сообщение - команда, то на этом наши приключения закончились.
		if res {
			return
		}

		// Здесь же, если фраза обращена к боту (ник или имя), выставляем флажок, что надо ответить
		quotedBotUsername := regexp.QuoteMeta(me.Result.Username)
		quotedBotFirstName := regexp.QuoteMeta(me.Result.FirstName)
		quotedBotLastName := regexp.QuoteMeta(me.Result.LastName)

		if messageContainsBotUsername, err := regexp.MatchString(quotedBotUsername, msg.Message.Text); err != nil {
			log.Errorf("An error occured while matching message with bot username: %s", err)

			return
		} else if messageContainsBotUsername {
			rmsg.Misc.Answer = 1
			m := regexp.MustCompile(quotedBotUsername)
			msg.Message.Text = m.ReplaceAllString(msg.Message.Text, "")
		}

		// TODO: обращение может содержать знаки препинания, пробелы перед, пробелы после имени бота. Надо бы это поддержать.
		// TODO: унести вырезание имени бота в отдельную процедуру.
		switch {
		// Если у нас есть и имя и фамилия бота, то попробуем их слепить вместе, через "1 и более" пробел и вырезать
		// полученное из входящей фразы.
		case quotedBotFirstName != "" && quotedBotLastName != "":
			messageContainsBotName, err := regexp.MatchString(
				quotedBotFirstName+"[[:space:]]+"+quotedBotLastName,
				msg.Message.Text,
			)

			if err != nil {
				log.Errorf("An error occured while matching message with bot firstname and lastname: %s", err)

				return
			}

			if messageContainsBotName {
				rmsg.Misc.Answer = 1
				m := regexp.MustCompile(quotedBotFirstName + "[[:space:]]+" + quotedBotLastName)
				msg.Message.Text = m.ReplaceAllString(msg.Message.Text, "")
			}
		// Если у нас есть только имя, то попробуем его вырезать из входящей фразы.
		case quotedBotFirstName != "":
			messageContainsBotName, err := regexp.MatchString(quotedBotFirstName, msg.Message.Text)

			if err != nil {
				log.Errorf("An error occured while matching message with bot firstname: %s", err)

				return
			}

			if messageContainsBotName {
				rmsg.Misc.Answer = 1
				m := regexp.MustCompile(quotedBotFirstName)
				msg.Message.Text = m.ReplaceAllString(msg.Message.Text, "")
			}
		// Если у нас есть только фамилия, то попробуем его вырезать из входящей фразы.
		case quotedBotLastName != "":
			messageContainsBotName, err := regexp.MatchString(quotedBotLastName, msg.Message.Text)

			if err != nil {
				log.Errorf("An error occured while matching message with bot lastname: %s", err)

				return
			}

			if messageContainsBotName {
				rmsg.Misc.Answer = 1
				m := regexp.MustCompile(quotedBotLastName)
				msg.Message.Text = m.ReplaceAllString(msg.Message.Text, "")
			}
		}

		// Засылаем фразу в misc-канал (в роутер).
		rmsg.Chatid = fmt.Sprintf("%d", msg.Message.Chat.ID)
		rmsg.Userid = fmt.Sprintf("%d", msg.Message.From.ID)
		rmsg.Message = msg.Message.Text
		rmsg.Mode = "public"
		rmsg.Plugin = "telegram"
		rmsg.From = "telegram"
		rmsg.Misc.Csign = Config.Csign
		rmsg.Misc.Fwdcnt = 1

		rmsg.Misc.Botnick = ConstructPartialUserUsername(me.Result)
		rmsg.Misc.Username = ConstructPartialUserUsername(msg.Message.From)

		rmsg.Misc.GoodMorning = 0
		// Форматирование нужно только для вывода некоторых ответов на команды, команды мы ловим выше по тексту, так что
		// смело ставим тут 0.
		rmsg.Misc.Msgformat = 0

		data, err := json.Marshal(rmsg)

		if err != nil {
			log.Warnf("Unable to to serialize message for redis: %s", err)

			return
		}

		// Заталкиваем наш json в редиску.
		if err := RedisClient.Publish(ctx, Config.Redis.Channel, data).Err(); err != nil {
			log.Warnf("Unable to send data to redis channel %s: %s", Config.Redis.Channel, err)
		} else {
			log.Debugf("Sent msg to redis channel %s: %s", Config.Redis.Channel, string(data))
		}

		// case "channel":
		// Тут мы ничего сделать не можем.
	} //nolint:wsl
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
