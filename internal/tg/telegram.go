package tg

import (
	"fmt"
	"strconv"

	"aleesa-telegram-go/internal/log"

	"github.com/NicoNex/echotron/v3"
)

// Telega основная горутинка, реализующая бота.
func Telega(c MyConfig) {
	tg = echotron.NewAPI(c.Telegram.Token)

	for u := range echotron.PollingUpdates(c.Telegram.Token) {
		telegramMsgParser(u)
	}
}

func delMsg(msg *echotron.Update) {
	if res, err := tg.DeleteMessage(msg.ChatID(), msg.Message.ID); err != nil {
		chat := ConstructFullChatName(msg.Message.SenderChat)
		user := ConstructFullUserName(msg.Message.From)

		log.Errorf("Unable to delete message %d in chat %s from %s via telegram api: %s", msg.Message.ID, chat, user, err)
	} else if !res.Ok {
		chat := ConstructFullChatName(msg.Message.SenderChat)
		user := ConstructFullUserName(msg.Message.From)

		log.Errorf("Unable to delete message %d in chat %s from %s via telegram api: %s", msg.ID, chat, user, res.Description)
	}
}

// ConstructFullUserName выковыривает из сообщения полный username, в формате @username FirstName LastName (id).
func ConstructFullUserName(u *echotron.User) string {
	user := fmt.Sprintf("(%d)", u.ID)

	if u.LastName != "" {
		user = fmt.Sprintf("%s %s", u.LastName, user)
	}

	if u.FirstName != "" {
		user = fmt.Sprintf("%s %s", u.FirstName, user)
	}

	if u.Username != "" {
		user = fmt.Sprintf("@%s %s", u.Username, user)
	}

	return user
}

// ConstructFullChatName выковыривает из сообщения полный username чата, в формате @username FirstName LastName (id).
func ConstructFullChatName(c *echotron.Chat) string {
	chat := fmt.Sprintf("(%d)", c.ID)

	if c.LastName != "" {
		chat = fmt.Sprintf("%s %s", c.LastName, chat)
	}

	if c.FirstName != "" {
		chat = fmt.Sprintf("%s %s", c.FirstName, chat)
	}

	if c.Username != "" {
		chat = fmt.Sprintf("@%s %s", c.Username, chat)
	}

	return chat
}

// ConstructPartialUserUsername пытается найти и вытащить username, если такового нет, вытаскивает First/Last Name, если
// такового нет, то возвращает ID.
func ConstructPartialUserUsername(u *echotron.User) string {
	switch {
	case u.Username != "":
		return "@" + u.Username

	case u.FirstName != "" && u.LastName != "":
		return u.FirstName + " " + u.LastName

	case u.FirstName != "":
		return u.FirstName

	case u.LastName != "":
		return u.LastName

	default:
		return strconv.FormatInt(u.ID, 10)
	}
}

// ConstructPartialChatUsername пытается найти и вытащить username, если такового нет, вытаскивает First/Last Name, если
// такового нет, то возвращает ID.
func ConstructPartialChatUsername(c *echotron.Chat) string {
	switch {
	case c.Username != "":
		return "@" + c.Username

	case c.FirstName != "" && c.LastName != "":
		return c.FirstName + " " + c.LastName

	case c.FirstName != "":
		return c.FirstName

	case c.LastName != "":
		return c.LastName

	default:
		return strconv.FormatInt(c.ID, 10)
	}
}

// ConstructUserFirstLastName Пытается найти и вытащить first name и last name пользователя, если не получается, то
// вначале пытается фоллбэчится на first name, потом на last name, потом на username.
func ConstructUserFirstLastName(u *echotron.User) string {
	var user string

	switch {
	case u.FirstName != "" && u.LastName != "":
		user = u.FirstName + " " + u.LastName
	case u.FirstName != "":
		user = u.FirstName
	case u.LastName != "":
		user = u.LastName
	case u.Username != "":
		user = "@" + u.Username
	default:
		user = strconv.FormatInt(u.ID, 10)
	}

	return user
}

// ConstructTelegramHighlightName генерирует имя пользователя, которое триггерит на стороне клиента ивент меншена
// определённого пользователя.
func ConstructTelegramHighlightName(u *echotron.User) string {
	var (
		username string
		link     string
	)

	link = fmt.Sprintf("tg://user?id=%d", u.ID)

	// Тут мы предполагаем, что как минимум либо firstname либо lastname либо username всегда есть. По сути так оно и
	// должно быть.
	switch {
	case u.FirstName != "" && u.LastName != "":
		username = fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	case u.FirstName != "":
		username = u.FirstName
	case u.LastName != "":
		username = u.LastName
	default:
		username = u.Username
	}

	return fmt.Sprintf("[%s](%s)", username, link)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
