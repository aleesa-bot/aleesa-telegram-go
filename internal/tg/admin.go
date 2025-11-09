// Package tg implements most of aleesa-telegram-go functionality.
package tg

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/NicoNex/echotron/v3"
	"github.com/davecgh/go-spew/spew"
)

// Admin парсит команды семейства admin. Возвращает true, если команда была опознана как admin.
func Admin(msg *echotron.Update) (bool, error) {
	var (
		c      = msg.Message.Text[len(Config.Csign):]
		answer string
	)

	switch {
	case c == "admin" || c == "админ":
		answer = "```\n"

		answer += Config.Csign + "admin censor          - показать список состояния типов сообщений\n"
		answer += Config.Csign + "админ ценз            - показать список состояния типов сообщений\n"
		answer += Config.Csign + "admin censor type #   - где 1 - вкл, 0 - выкл цензуры для означенного типа сообщений\n"
		answer += Config.Csign + "админ ценз тип #      - где 1 - вкл, 0 - выкл цензуры для означенного типа сообщений\n"
		answer += Config.Csign + "admin fortune         - показываем ли с утра фортунку для чата\n"
		answer += Config.Csign + "admin фортунка        - показываем ли с утра фортунку для чата\n"
		answer += Config.Csign + "admin fortune #       - где 1 - вкл, 0 - выкл фортунку с утра\n"
		answer += Config.Csign + "admin фортунка #      - где 1 - вкл, 0 - выкл фортунку с утра\n"
		answer += Config.Csign + "admin greet           - приветствуем ли новых участников чата\n"
		answer += Config.Csign + "admin приветствие     - приветствуем ли новых участников чата\n"
		answer += Config.Csign + "admin greet #         - где 1 - вкл, 0 - выкл приветствия новых участников чата\n"
		answer += Config.Csign + "admin приветствие #   - где 1 - вкл, 0 - выкл приветствия новых участников чата\n"
		answer += Config.Csign + "admin goodbye         - прощаемся ли с ушедшими участниками чата\n"
		answer += Config.Csign + "admin прощание        - прощаемся ли с ушедшими участниками чата\n"
		answer += Config.Csign + "admin goodbye #       - где 1 - вкл, 0 - выкл прощания с ушедшими участниками чата\n"
		answer += Config.Csign + "admin прощание #      - где 1 - вкл, 0 - выкл прощания с ушедшими участниками чата\n"
		answer += Config.Csign + "admin oboobs #        - где 1 - вкл, 0 - выкл плагина oboobs\n"

		answer += fmt.Sprintf(
			"%sadmin oboobs          - показываем ли сисечки по просьбе участников чата (команды %stits, %stities, %sboobs, %sboobies, %sсиси, %sсисечки)\n",
			Config.Csign, Config.Csign, Config.Csign, Config.Csign, Config.Csign, Config.Csign, Config.Csign,
		)

		answer += Config.Csign + "admin obutts #        - где 1 - вкл, 0 - выкл плагина obutts\n"

		answer += fmt.Sprintf(
			"%sadmin obutts          - показываем ли попки по просьбе участников чата (команды %sass, %sbutt, %sbooty, %sпопа, %sпопка)\n",
			Config.Csign, Config.Csign, Config.Csign, Config.Csign, Config.Csign, Config.Csign,
		)

		answer += Config.Csign + "admin chan_msg        - оставляем ли сообщения присланные от имени (других) каналов\n"
		answer += Config.Csign + "admin chan_msg #      - где 1 - оставляем, 0 - удаляем\n"
		answer += Config.Csign + "admin ban userid sec  - выдаём ban указанному user-у на указанное количество секунд (от 30 сек до 1 года), доступно только создателю чата\n"
		answer += Config.Csign + "admin mute userid sec - выдаём mute указанному user-у на указанное количество секунд (от 30 сек до 1 года), доступно только создателю чата\n"
		answer += Config.Csign + "admin admin mute      - разрешено ли обычным админам мьютить участников чата через бота (если бот - админ), (создатель чата всегда может попросить бота-админа замьютить обычного участника чата)\n"

		answer += Config.Csign + "%sadmin admin mute #    - где 1 - разрешено, 0 - не разрешено\n"
		answer += "\n"
		answer += "Типы сообщений:\naudio voice photo video video_note animation sticker dice game poll document\n"
		answer += "```"

	case regexp.MustCompile("^(admin|админ)[[:space:]](censor|ценз)$").MatchString(c):
		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "VoiceMsg") == "1" {
			answer += "Тип сообщения voice удаляется\n"
		} else {
			answer += "Тип сообщения voice не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "AudioMsg") == "1" {
			answer += "Тип сообщения audio удаляется\n"
		} else {
			answer += "Тип сообщения audio не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "PhotoMsg") == "1" {
			answer += "Тип сообщения photo удаляется\n"
		} else {
			answer += "Тип сообщения photo не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "VideoMsg") == "1" {
			answer += "Тип сообщения video удаляется\n"
		} else {
			answer += "Тип сообщения video не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "VideoNoteMsg") == "1" {
			answer += "Тип сообщения video_note удаляется\n"
		} else {
			answer += "Тип сообщения video_note не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "AnimationMsg") == "1" {
			answer += "Тип сообщения animation удаляется\n"
		} else {
			answer += "Тип сообщения animation не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "StickerMsg") == "1" {
			answer += "Тип сообщения sticker удаляется\n"
		} else {
			answer += "Тип сообщения sticker не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "DiceMsg") == "1" {
			answer += "Тип сообщения dice удаляется\n"
		} else {
			answer += "Тип сообщения dice не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "GameMsg") == "1" {
			answer += "Тип сообщения game удаляется\n"
		} else {
			answer += "Тип сообщения game не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "PollMsg") == "1" {
			answer += "Тип сообщения poll удаляется\n"
		} else {
			answer += "Тип сообщения poll не удаляется\n"
		}

		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "DocumentMsg") == "1" {
			answer += "Тип сообщения document удаляется\n"
		} else {
			answer += "Тип сообщения document не удаляется\n"
		}

	case regexp.MustCompile("^(admin|админ)[[:space:]](censor|ценз)[[:space:]].+$").MatchString(c):
		// Снова регулярка, это не оч эффективно, но зато просто.
		// Выхватываем тип сообщения и согласно типу ставим настройку.
		cmdArray := regexp.MustCompile("[[:space:]]").Split(c, 4)

		switch cmdArray[2] {
		case "voice":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "VoiceMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с voice будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "VoiceMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с voice будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "audio":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "AudioMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с audio будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "AudioMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с audio будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "photo":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "PhotoMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с photo будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "PhotoMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с photo будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "video":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "VideoMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с video будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "VideoMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с video будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "video_note":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "VideoNoteMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с video_note будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "VideoNoteMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с video_note будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "animation":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "AnimationMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с animation будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "AnimationMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с animation будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "sticker":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "StickerMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с sticker будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "StickerMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с sticker будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "dice":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "DiceMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с dice будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "diceMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с dice будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "game":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "GameMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с game будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "GameMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с game будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "poll":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "PollMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с poll будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "PollMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с poll будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		case "document":
			switch cmdArray[3] {
			case "1":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "DocumentMsg", "1"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с document будут оставаться"

			case "0":
				if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "DocumentMsg", "0"); err != nil {
					return true, err
				}

				answer += "Теперь сообщения с document будут удаляться"

			default:
				answer += "Не понимаю о чём ты."
			}

		default:
			answer += "Не понимаю о чём ты."
		}

	case regexp.MustCompile("^(admin|админ)[[:space:]](fortune|фортунка)$").MatchString(c):
		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "FortuneMsg") == "1" {
			answer += "Фортунки с утра показываются"
		} else {
			answer += "Фортунки с утра не показываются"
		}

	case regexp.MustCompile("^(admin|админ)[[:space:]](fortune|фортунка)[[:space:]].*$").MatchString(c):
		// Снова регулярка, это не оч эффективно, но зато просто.
		cmdArray := regexp.MustCompile("[[:space:]]").Split(c, 3)

		// TODO: На самом деле надо сохранять список чятиков, в которые надо отправлять сообщения. Из этой базы
		//       мы не сможем вытащить список чятиков.
		switch cmdArray[2] {
		case "1":
			if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "FortuneMsg", "1"); err != nil {
				return true, err
			}

			answer += "Теперь сообщения с фортункой будут отправляться с утра"

		case "0":
			if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "FortuneMsg", "0"); err != nil {
				return true, err
			}

			answer += "Теперь сообщения с фортункой не будут отправляться с утра"

		default:
			answer += "Не понимаю о чём ты."
		}

	case regexp.MustCompile("^(admin|админ)[[:space:]](greet|приветствие)$").MatchString(c):
		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "GreetMsg") == "1" {
			answer += "Приветствуем вновьприбывших участников чатика"
		} else {
			answer += "Не приветствуем вновьприбывших участников чатика"
		}

	case regexp.MustCompile("^(admin|админ)[[:space:]](greet|приветствие)[[:space:]].*$").MatchString(c):
		// Снова регулярка, это не оч эффективно, но зато просто.
		cmdArray := regexp.MustCompile("[[:space:]]").Split(c, 3)

		// TODO: На самом деле надо сохранять список чятиков, в которые надо отправлять сообщения. Из этой базы
		//       мы не сможем вытащить список чятиков.
		switch cmdArray[2] {
		case "1":
			if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "GreetMsg", "1"); err != nil {
				return true, err
			}

			answer += "Будем приветствовать вновьприбывших участников чата"

		case "0":
			if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "GreetMsg", "0"); err != nil {
				return true, err
			}

			answer += "Не будем приветствовать вновьприбывших участников чата"

		default:
			answer += "Не понимаю, о чём ты."
		}

	case regexp.MustCompile("^(admin|админ)[[:space:]](goodbye|прощание)$").MatchString(c):
		if GetSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "GoodbyeMsg") == "1" {
			answer += "Прощаемся с ушедшими участниками чатика"
		} else {
			answer += "Не прощаемся с ушедшими участниками чатика"
		}

	case regexp.MustCompile("^(admin|админ)[[:space:]](goodbye|прощание)[[:space:]].*$").MatchString(c):
		// Снова регулярка, это не оч эффективно, но зато просто.
		cmdArray := regexp.MustCompile("[[:space:]]").Split(c, 3)

		switch cmdArray[2] {
		case "1":
			if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "GoodbyeMsg", "1"); err != nil {
				return true, err
			}

			answer += "Будем прощаться с ушедшими участниками чата"

		case "0":
			if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "GoodbyeMsg", "0"); err != nil {
				return true, err
			}

			answer += "Не будем прощаться с ушедшими участниками чата"

		default:
			answer += "Не понимаю, о чём ты."
		}

	case regexp.MustCompile("^(admin|админ)[[:space:]](chan_msg)[[:space:]].*$").MatchString(c):
		// Снова регулярка, это не оч эффективно, но зато просто.
		cmdArray := regexp.MustCompile("[[:space:]]").Split(c, 3)

		switch cmdArray[2] {
		case "1":
			if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "ChanMsg", "1"); err != nil {
				return true, err
			}

			answer += "Будем удалять сообщения, написанные от имени канала"

		case "0":
			if err := SaveSetting(strconv.FormatInt(msg.Message.Chat.ID, 10), "ChanMsg", "0"); err != nil {
				return true, err
			}

			answer += "Не будем удалять сообщения, написанные от имени канала"

		default:
			answer += "Не понимаю, о чём ты."
		}
	}

	resp, err := tg.SendMessage(
		answer,
		msg.Message.Chat.ID,
		&echotron.MessageOptions{ParseMode: "MarkdownV2"},
	)

	if err != nil || !resp.Ok {
		// TODO: поддержать миграцию группы в супергруппу, поддержать вариант, когда бот замьючен.
		// Красиво оформить ошибку, с полями итд, как tracedump, только ошибка.
		// N.B. тут может быть сообщение о том, что группа превратилась в супергруппу, или что бот не имеет прав писать
		// сообщения в чятик. Это надо бы отслеживать и хэндлить.
		return true, fmt.Errorf(
			"unable to send message to telegram api: %w\nResponse dump: %s",
			err,
			spew.Sdump(resp),
		)
	}

	return true, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
