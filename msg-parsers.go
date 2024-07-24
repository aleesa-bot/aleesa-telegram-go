package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/NicoNex/echotron/v3"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// telego msg parser

// redisMsgParser парсит json-чики прилетевшие из REDIS-ки, причём json-чики должны быть относительно валидными.
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

	// Сообщение о том, что этот чятик изменился, например, превратился в супергруппу.
	if (msg.Message.MigrateFromChatID < 0) && (msg.Message.MigrateToChatID < 0) { //nolint: revive,staticcheck
		// TODO: поддержать миграцию настроек чата на новый chatId.
		// TODO: подумать, что можно сделать с настройками бэкэндов. Кажись, ничего, но надо глянуть.
	}

	// Люди пришли в чят.
	if msg.Message.NewChatMembers != nil { //nolint: revive,staticcheck
		// TODO: проверить, надо ли приветствовать, выбрать одну из приветственных фраз, потянуть время, типа, мы пишем
		// текст и поприветствовать вновь прибывшего.
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
	}

	// Если сообщение было зацензурено, то дальше его обрабатывать не надо.
	if Censor(msg) {
		return
	}

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
		rmsg.Misc.Csign = config.Csign
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
			if err := redisClient.Publish(ctx, config.Redis.Channel, data).Err(); err != nil {
				log.Warnf("Unable to send data to redis channel %s: %s", config.Redis.Channel, err)
			} else {
				log.Debugf("Sent msg to redis channel %s: %s", config.Redis.Channel, string(data))
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

		// TODO: обращение может содержать знаки препинания, пробелы перед, пробелы после. Надо бы это поддержать.
		switch {
		// Если у на есть и имя и фамилия бота, то попробуем их слепить вместе, через "1 и более" пробел и вырезать
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
		rmsg.Misc.Csign = config.Csign
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
		if err := redisClient.Publish(ctx, config.Redis.Channel, data).Err(); err != nil {
			log.Warnf("Unable to send data to redis channel %s: %s", config.Redis.Channel, err)
		} else {
			log.Debugf("Sent msg to redis channel %s: %s", config.Redis.Channel, string(data))
		}

		// case "channel":
		// Тут мы ничего сделать не можем.
	} //nolint:wsl
}

// cmdParser разбирает сообщение как команду и засылает его в роутер, попутно обкладывая корректными значениями
// параметров. Возвращает true если это была команда.
func cmdParser(me echotron.APIResponseUser, cmd *echotron.Update) (bool, error) {
	var (
		err   error
		smsg  sMsg
		bingo bool
	)

	smsg.From = "telegram"
	smsg.Plugin = "telegram"
	smsg.Chatid = fmt.Sprintf("%d", cmd.Message.Chat.ID)
	smsg.Userid = fmt.Sprintf("%d", cmd.Message.From.ID)
	smsg.Message = cmd.Message.Text
	smsg.Misc.Answer = 1

	// Тоже самое можно проставить выковырнув значение из cmd.Message.Chat.Type, где group/supergroup/channel - это
	// public, а остальное - private, но по id проще/быстрее, хоть это и хак.
	if cmd.Message.Chat.ID >= 0 {
		smsg.Mode = "private"
	} else {
		smsg.Mode = "public"
	}

	smsg.Misc.Botnick = ConstructPartialUserUsername(me.Result)

	// Предполагаем, что длина cmd.Message.Text всегда больше или равна длине config.Csign.
	if cmd.Message.Text[0:len(config.Csign)] == config.Csign {
		// Повторно проверяем, что текст является простой командой.
		command := cmd.Message.Text[len(config.Csign):]

		// TODO: поддержать команды admin и admin *.
		if command == "помощь" || command == "help" {
			help := "```\n"
			help += fmt.Sprintf("%shelp | %sпомощь             - список команд", config.Csign, config.Csign)
			help += fmt.Sprintf("%sanek | %sанек | %sанекдот    - рандомный анекдот с anekdot.ru", config.Csign, config.Csign, config.Csign)
			help += fmt.Sprintf("%sbuni                       - рандомный стрип hapi buni", config.Csign)
			help += fmt.Sprintf("%sbunny | %srabbit | %sкролик  - кролик", config.Csign, config.Csign, config.Csign)
			help += fmt.Sprintf("%scat | %sкис                 - кошечка", config.Csign, config.Csign)
			help += fmt.Sprintf("%scoin | %sмонетка            - подбросить монетку - орёл или решка?", config.Csign, config.Csign)
			help += fmt.Sprintf("%sdig | %sкопать              - заняться археологией", config.Csign, config.Csign)
			help += fmt.Sprintf("%sdrink | %sпраздник          - какой сегодня праздник?", config.Csign, config.Csign)
			help += fmt.Sprintf("%sfish | %sрыба | %sрыбка      - порыбачить", config.Csign, config.Csign, config.Csign)
			help += fmt.Sprintf("%sfishing | %sрыбалка         - порыбачить", config.Csign, config.Csign)
			help += fmt.Sprintf("%sf | %sф                     - рандомная фраза из сборника цитат fortune_mod", config.Csign, config.Csign)
			help += fmt.Sprintf("%sfortune | %sфортунка        - рандомная фраза из сборника цитат fortune_mod", config.Csign, config.Csign)
			help += fmt.Sprintf("%sfox | %sлис                 - лисичка", config.Csign, config.Csign)
			help += fmt.Sprintf("%sfriday | %sпятница          - а не пятница ли сегодня?", config.Csign, config.Csign)
			help += fmt.Sprintf("%sfrog | %sлягушка            - лягушка", config.Csign, config.Csign)
			help += fmt.Sprintf("%shorse | %sлошадка           - лошадка", config.Csign, config.Csign)
			help += fmt.Sprintf("%slat | %sлат                 - сгенерить фразу из крылатых латинских выражений", config.Csign, config.Csign)
			help += fmt.Sprintf("%smonkeyuser                 - рандомный стрип MonkeyUser", config.Csign)
			help += fmt.Sprintf("%sowl | %sсова | %sсыч         - сова", config.Csign, config.Csign, config.Csign)
			help += fmt.Sprintf("%sproverb | %sпословица       - рандомная русская пословица", config.Csign, config.Csign)
			help += fmt.Sprintf("%sping | %sпинг               - попинговать бота", config.Csign, config.Csign)
			help += fmt.Sprintf("%sroll | %sdice | %sкости      - бросить кости", config.Csign, config.Csign, config.Csign)
			help += fmt.Sprintf("%ssnail | %sулитка            - улитка", config.Csign, config.Csign)
			help += fmt.Sprintf("%sver | %sversion | %sверсия   - что-то про версию ПО", config.Csign, config.Csign, config.Csign)
			help += fmt.Sprintf("%sw город | %sп город         - погода в указанном городе", config.Csign, config.Csign)
			help += fmt.Sprintf("%sweather город              - погода в указанном городе", config.Csign)
			help += fmt.Sprintf("%sпогода город               - погода в указанном городе", config.Csign)
			help += fmt.Sprintf("%sпогодка город              - погода в указанном городе", config.Csign)
			help += fmt.Sprintf("%sпогадка город              - погода в указанном городе", config.Csign)
			help += fmt.Sprintf("%sxkcd                       - рандомный стрип с сайта xkcd.ru", config.Csign)
			help += fmt.Sprintf("%skarma фраза | %sкарма фраза - посмотреть карму фразы", config.Csign, config.Csign)
			help += "```\n"
			help += "Но на самом деле я бот больше для общения, чем для исполнения команд.\n"
			help += "Поговоришь со мной?"

			var opts *echotron.MessageOptions

			chatid := cmd.Message.Chat.ID

			resp, err := tg.SendMessage(help, chatid, opts)

			if err != nil || !resp.Ok {
				// TODO: поддержать миграцию группы в супергруппу, поддержать вариант, когда бот замьючен.
				// Красиво оформить ошибку, с полями итд, как tracedump, только ошибка.
				// N.B. тут может быть сообщение о том, что группа превратилась в супергруппу, или что бот не имеет прав писать
				// сообщения в чятик. Это надо бы отслеживать и хэндлить.
				log.Errorf("Unable to send message to telegram api: %s", err)
				log.Errorf("Response dump: %s", spew.Sdump(resp))
			}
		}

		// Команды в одно слово.
		cmds := []string{
			"ping", "пинг", "пинх", "pong", "понг", "понх", "coin", "монетка", "roll", "dice", "кости", "ver",
			"version", "версия", "хэлп", "halp", "kde", "кде", "lat", "лат", "friday", "пятница", "proverb",
			"пословица", "пословиться", "fortune", "фортунка", "f", "ф", "anek", "анек", "анекдот", "buni", "cat",
			"кис", "drink", "праздник", "fox", "лис", "frog", "лягушка", "horse", "лошадь", "лошадка", "monkeyuser",
			"owl", "сова", "сыч", "rabbit", "bunny", "кролик", "snail", "улитка", "tits", "boobs", "tities", "boobies",
			"сиси", "сисечки", "butt", "booty", "ass", "попа", "попка", "xkcd", "dig", "копать", "fish", "fishing",
			"рыба", "рыбка", "рыбалка", "karma", "карма", "fuck",
		}

		for _, c := range cmds {
			if command == c {
				bingo = true

				break
			}
		}

		// Не нашлось команды в одно слово. Поищем команды с одним аргументом.
		if !bingo {
			cmds = []string{
				"w", "п", "weather", "погода", "погодка", "погадка", "karma", "карма",
			}

			for _, c := range cmds {
				cRe := fmt.Sprintf("^%s ", c)
				r := regexp.MustCompile(cRe)

				if r.MatchString(command) {
					bingo = true

					break
				}
			}
		}

		// Что-то странное пришло. Залоггируем и ничего делать не будем. Просто свалим.
		if !bingo {
			log.Infof(
				"Strange command from %s %s: %s",
				ConstructFullChatName(&cmd.Message.Chat),
				ConstructFullUserName(cmd.Message.From),
				cmd.Message.Text,
			)

			return true, err
		}

		// Для некоторых команд надо подсвечивать имя пользователя в ответе.
		cmds = []string{
			"dig", "копать", "fish", "fishing", "рыба", "рыбка", "рыбалка",
		}

		for _, c := range cmds {
			// В vscode возникает ошибка при использовании if-а.
			switch {
			case command == c:
				smsg.Misc.Username = ConstructTelegramHighlightName(cmd.Message.From)
				smsg.Misc.Msgformat = 1

				break //nolint: gosimple
			}
		}

		// Если у нас супергруппа с тредами, надо проставить thread id.
		if cmd.Message.Chat.Type == "supergroup" && cmd.Message.IsTopicMessage && cmd.Message.ThreadID != 0 {
			smsg.ThreadID = fmt.Sprintf("%d", cmd.Message.ThreadID)
		}

		// Отправляем сообщение в редиску.
		data, err := json.Marshal(smsg)

		if err != nil {
			log.Warnf("Unable to to serialize message for redis: %s", err)

			return true, err
		}

		// Заталкиваем наш json в редиску.
		if err := redisClient.Publish(ctx, config.Redis.Channel, data).Err(); err != nil {
			log.Warnf("Unable to send data to redis channel %s: %s", config.Redis.Channel, err)
		} else {
			log.Debugf("Sent msg to redis channel %s: %s", config.Redis.Channel, string(data))
		}

		return true, err
	}

	return false, err
}

// Censor парсит сообщения в поисках непотребных данных и если он их находит, то сообщение удаляется.
// Непотребными могут быть аудиосообщения, аудиофайлы, видеосообщения, сообщения от имени других каналов итп.
// Это могут настроить админы чятика через команду !admin censor.
func Censor(msg *echotron.Update) bool {
	result := false

	switch {
	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "VoiceMsg") == "1":
		result = true

		// Предполагаем, что у voice-ов здесь всегда не ноль.
		if msg.Message.Voice.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "AudioMsg") == "1":
		// Предполагаем, что у аудио здесь всегда не ноль.
		if msg.Message.Audio.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "PhotoMsg") == "1":
		// Обычное сообщение не содержит фоток.
		if len(msg.Message.Photo) != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "VideoMsg") == "1":
		// Предполагаем, что у видео здесь всегда не ноль.
		if msg.Message.Video.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "VideoNoteMsg") == "1":
		// Предполагаем, что у видео-заметки здесь всегда не ноль.
		if msg.Message.VideoNote.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "AnimationMsg") == "1":
		// Предполагаем, что у анимации здесь всегда не ноль.
		if msg.Message.Animation.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "StickerMsg") == "1":
		// Предполагаем, что FileID не пустое только у стикера.
		if msg.Message.Sticker.FileID != "" {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "DiceMsg") == "1":
		// Предполагаем, что Value > 0 только у дайса.
		if msg.Message.Dice.Value != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "GameMsg") == "1":
		// Предполагаем, что title только у game-а.
		if msg.Message.Game.Title != "" {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "PollMsg") == "1":
		// Предполагаем, что title только у game-а.
		if msg.Message.Poll.Question != "" {
			delMsg(msg)

			result = true
		}

	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "DocumentMsg") == "1":
		// Предполагаем, что FileID есть только у document-а.
		if msg.Message.Document.FileID != "" {
			delMsg(msg)

			result = true
		}

	// Некоторые рекламные товарищи пытаются срать своими каналами в чятик это тоже можно зацензурить ботом и это
	// пидорство он будет удалять asap.
	// 136817688 - это специальный id пользователя, который принимает облик канала, на него можно нажать и попасть
	//             на рекламируемый канал.
	case GetSetting(fmt.Sprintf("%d", msg.ChatID()), "ChanMsg") == "1":
		if msg.Message.From.ID == 136817688 {
			delMsg(msg)

			result = true
		}
	}

	return result
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
