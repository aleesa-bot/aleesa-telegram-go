package tg

import (
	"strconv"

	"github.com/NicoNex/echotron/v3"
)

// Censor парсит сообщения в поисках непотребных данных и если он их находит, то сообщение удаляется.
// Непотребными могут быть аудиосообщения, аудиофайлы, видеосообщения, сообщения от имени других каналов итп.
// Это могут настроить админы чятика через команду !admin censor.
func Censor(msg *echotron.Update) bool {
	var (
		result bool
		chatID = msg.ChatID()
	)

	switch {
	case GetSetting(strconv.FormatInt(chatID, 10), "VoiceMsg") == "1":
		result = true

		// Предполагаем, что у voice-ов здесь всегда не ноль.
		if msg.Message.Voice.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "AudioMsg") == "1":
		// Предполагаем, что у аудио здесь всегда не ноль.
		if msg.Message.Audio.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "PhotoMsg") == "1":
		// Обычное сообщение не содержит фоток.
		if len(msg.Message.Photo) != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "VideoMsg") == "1":
		// Предполагаем, что у видео здесь всегда не ноль.
		if msg.Message.Video.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "VideoNoteMsg") == "1":
		// Предполагаем, что у видео-заметки здесь всегда не ноль.
		if msg.Message.VideoNote.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "AnimationMsg") == "1":
		// Предполагаем, что у анимации здесь всегда не ноль.
		if msg.Message.Animation.Duration != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "StickerMsg") == "1":
		// Предполагаем, что FileID не пустое только у стикера.
		if msg.Message.Sticker.FileID != "" {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "DiceMsg") == "1":
		// Предполагаем, что Value > 0 только у дайса.
		if msg.Message.Dice.Value != 0 {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "GameMsg") == "1":
		// Предполагаем, что title только у game-а.
		if msg.Message.Game.Title != "" {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "PollMsg") == "1":
		// Предполагаем, что title только у game-а.
		if msg.Message.Poll.Question != "" {
			delMsg(msg)

			result = true
		}

	case GetSetting(strconv.FormatInt(chatID, 10), "DocumentMsg") == "1":
		// Предполагаем, что FileID есть только у document-а.
		if msg.Message.Document.FileID != "" {
			delMsg(msg)

			result = true
		}

	// Некоторые рекламные товарищи пытаются срать своими каналами в чятик это тоже можно зацензурить ботом и это
	// пидорство он будет удалять asap.
	// 136817688 - это специальный id пользователя, который принимает облик канала, на него можно нажать и попасть
	//             на рекламируемый канал.
	case GetSetting(strconv.FormatInt(chatID, 10), "ChanMsg") == "1":
		if msg.Message.From.ID == 136817688 {
			delMsg(msg)

			result = true
		}
	}

	return result
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
