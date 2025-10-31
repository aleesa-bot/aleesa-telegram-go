package tg

import (
	"aleesa-telegram-go/internal/log"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/NicoNex/echotron/v3"
)

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

	// Предполагаем, что длина cmd.Message.Text всегда больше или равна длине Config.Csign.
	if cmd.Message.Text[0:len(Config.Csign)] == Config.Csign {
		// Повторно проверяем, что текст является простой командой.
		command := cmd.Message.Text[len(Config.Csign):]

		switch {
		case command == "помощь" || command == "help":
			return Help(cmd)

		// Хэндлер команд admin и admin *.
		case regexp.MustCompile("^(admin|админ)(.+)?$").MatchString(command):
			return Admin(cmd)
		}

		// Команды в одно слово.
		cmds := []string{
			"ping", "пинг", "пинх", "pong", "понг", "понх", "coin", "монетка", "roll", "dice", "кости", "ver",
			"version", "версия", "хэлп", "halp", "kde", "кде", "lat", "лат", "friday", "пятница", "proverb",
			"пословица", "пословиться", "fortune", "фортунка", "f", "ф", "anek", "анек", "анекдот", "buni", "cat",
			"кис", "drink", "праздник", "fox", "лис", "frog", "лягушка", "horse", "лошадь", "лошадка", "monkeyuser",
			"owl", "сова", "сыч", "rabbit", "bunny", "кролик", "snail", "улитка", "tits", "boobs", "tities", "boobies",
			"сиси", "сисечки", "butt", "booty", "ass", "попа", "попка", "xkcd", "dig", "копать", "fish", "fishing",
			"рыба", "рыбка", "рыбалка", "karma", "карма", "fuck", "weather", "погода", "w", "п", "погодка", "погадка",
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
		if err := RedisClient.Publish(ctx, Config.Redis.Channel, data).Err(); err != nil {
			log.Warnf("Unable to send data to redis channel %s: %s", Config.Redis.Channel, err)
		} else {
			log.Debugf("Sent msg to redis channel %s: %s", Config.Redis.Channel, string(data))
		}

		return true, err
	}

	return false, err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
