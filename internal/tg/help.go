package tg

import (
	"fmt"

	"github.com/NicoNex/echotron/v3"
	"github.com/davecgh/go-spew/spew"
)

// Help Выводит в чат сообщение с основными командами бота.
func Help(cmd *echotron.Update) (bool, error) {
	help := "```\n"

	help += fmt.Sprintf("%shelp | %sпомощь             - список команд", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sanek | %sанек | %sанекдот    - рандомный анекдот с anekdot.ru", Config.Csign, Config.Csign, Config.Csign)
	help += Config.Csign + "buni                       - рандомный стрип hapi buni"
	help += fmt.Sprintf("%sbunny | %srabbit | %sкролик  - кролик", Config.Csign, Config.Csign, Config.Csign)
	help += fmt.Sprintf("%scat | %sкис                 - кошечка", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%scoin | %sмонетка            - подбросить монетку - орёл или решка?", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sdig | %sкопать              - заняться археологией", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sdrink | %sпраздник          - какой сегодня праздник?", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sfish | %sрыба | %sрыбка      - порыбачить", Config.Csign, Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sfishing | %sрыбалка         - порыбачить", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sf | %sф                     - рандомная фраза из сборника цитат fortune_mod", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sfortune | %sфортунка        - рандомная фраза из сборника цитат fortune_mod", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sfox | %sлис                 - лисичка", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sfriday | %sпятница          - а не пятница ли сегодня?", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sfrog | %sлягушка            - лягушка", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%shorse | %sлошадка           - лошадка", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%slat | %sлат                 - сгенерить фразу из крылатых латинских выражений", Config.Csign, Config.Csign)
	help += Config.Csign + "monkeyuser                 - рандомный стрип MonkeyUser"
	help += fmt.Sprintf("%sowl | %sсова | %sсыч         - сова", Config.Csign, Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sproverb | %sпословица       - рандомная русская пословица", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sping | %sпинг               - попинговать бота", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sroll | %sdice | %sкости      - бросить кости", Config.Csign, Config.Csign, Config.Csign)
	help += fmt.Sprintf("%ssnail | %sулитка            - улитка", Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sver | %sversion | %sверсия   - что-то про версию ПО", Config.Csign, Config.Csign, Config.Csign)
	help += fmt.Sprintf("%sw <город> | %sп <город>     - погода в указанном городе", Config.Csign, Config.Csign)
	help += Config.Csign + "weather <город>            - погода в указанном городе"
	help += Config.Csign + "погода <город>             - погода в указанном городе"
	help += Config.Csign + "погодка <город>            - погода в указанном городе"
	help += Config.Csign + "погадка <город>            - погода в указанном городе"
	help += Config.Csign + "xkcd                       - рандомный стрип с сайта xkcd.ru"
	help += fmt.Sprintf("%skarma фраза | %sкарма фраза - посмотреть карму фразы", Config.Csign, Config.Csign)
	help += "```\n"
	help += "Но на самом деле я бот больше для общения, чем для исполнения команд.\n"
	help += "Поговоришь со мной?"

	resp, err := tg.SendMessage(
		help,
		cmd.Message.Chat.ID,
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
