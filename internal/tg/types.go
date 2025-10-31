package tg

type (
	// MyConfig структура, описывающая формат Конфига.
	MyConfig struct {
		Csign    string `json:"csign,omitempty"`
		DataDir  string `json:"data_dir,omitempty"`
		Log      string `json:"log,omitempty"`
		Loglevel string `json:"loglevel,omitempty"`

		Telegram struct {
			Token string `json:"token,omitempty"`
		} `json:"telegram"`

		Redis struct {
			Channel   string `json:"channel,omitempty"`
			MyChannel string `json:"my_channel,omitempty"`
			Server    string `json:"server,omitempty"`
			Port      int    `json:"port,omitempty"`
		} `json:"redis"`

		ForwardsMax int64 `json:"forwards_max,omitempty"`
	}

	// rMsg структура, описывающая формат входящего сообщения из pubsub-канала redis-ки.
	rMsg struct {
		From     string `json:"from,omitempty"`
		Chatid   string `json:"chatid,omitempty"`
		Userid   string `json:"userid,omitempty"`
		ThreadID string `json:"threadid,omitempty"`
		Message  string `json:"message,omitempty"`
		Plugin   string `json:"plugin,omitempty"`
		Mode     string `json:"mode,omitempty"`
		Misc     struct {
			Botnick string `json:"bot_nick,omitempty"`
			Csign   string `json:"csign,omitempty"`
			// Судя по перловой реализации, сюда всегда попадает [$fullname](tg://user?id=$userid).
			// Согласно докуметации, и перловой реализации бота, мы проставляем это поле тогда, когда мы точно знаем, что в
			// бэкэнде оно используется.
			Username    string `json:"username,omitempty"`
			Answer      int64  `json:"answer,omitempty"`
			Fwdcnt      int64  `json:"fwd_cnt,omitempty"`
			GoodMorning int64  `json:"good_morning,omitempty"`
			// В этом поле мы пишем 1, если бэкэнд должен в ответе форматировать сообщение в MarkdownV2, что какбэ
			// не оч хорошо, так как мы выносим особенности реализации бэкэнда сюда.
			Msgformat int64 `json:"msg_format,omitempty"`
		} `json:"Misc"`
	}

	// sMsg структура, описывающая исходящее сообщение в pubsub-канал redis-ки.
	sMsg struct {
		From     string `json:"from"`
		Chatid   string `json:"chatid"`
		Userid   string `json:"userid"`
		ThreadID string `json:"threadid,omitempty"`
		Message  string `json:"message"`
		Plugin   string `json:"plugin"`
		Mode     string `json:"mode"`
		Misc     struct {
			Botnick     string `json:"bot_nick"`
			Csign       string `json:"csign"`
			Username    string `json:"username"`
			Answer      int64  `json:"answer"`
			Fwdcnt      int64  `json:"fwd_cnt"`
			GoodMorning int64  `json:"good_morning"`
			Msgformat   int64  `json:"msg_format"`
		} `json:"misc"`
	}
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
