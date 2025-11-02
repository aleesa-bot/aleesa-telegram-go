package tg

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"slices"

	"aleesa-telegram-go/internal/log"

	"github.com/cockroachdb/errors/oserror"
	"github.com/cockroachdb/pebble"
)

// StoreKV сохраняет в указанной бд ключ и значение.
func StoreKV(db *pebble.DB, key string, value string) error {
	var (
		kArray = []byte(key)
		vArray = []byte(value)
	)

	err := db.Set(kArray, vArray, pebble.Sync)

	return err
}

// FetchV достаёт значение по ключу.
func FetchV(db *pebble.DB, key string) (string, error) {
	var (
		kArray = []byte(key)

		vArray      []byte
		valueString = ""
	)

	vArray, closer, err := db.Get(kArray)

	if err != nil {
		return valueString, err
	}

	valueString = string(vArray)
	err = closer.Close()

	return valueString, err
}

// InitChatListDB открывает для пользования базу со списком чатов.
func InitChatListDB() error {
	var (
		err     error
		options pebble.Options
		dbPath  = Config.DataDir + "/" + chatListDBName
	)

	// По дефолту ограничение ставится на мегабайты данных, а не на количество файлов, поэтому с дефолтными
	// настройками порождается огромное количество файлов. Умолчальное ограничение на количество файлов - 500 штук,
	// что нас не устраивает, поэтому немного снизим эту цифру до более приемлемых значений.
	options.L0CompactionFileThreshold = 8
	chatListDB, err = pebble.Open(dbPath, &options)

	if err != nil {
		return fmt.Errorf("Unable to open settings db, %s: %w\n", dbPath, err)
	}

	// Заполним слайс со списком чатов. Операция нудная.
	iterator, err := chatListDB.NewIter(&pebble.IterOptions{})

	if !iterator.First() {
		return errors.New("Unable to iterate slice with chats")
	}

	// Попробуем извлечь список ключей и заполнить слайс с чатами.
	for {
		// Закончили итерировать по списку.
		if !iterator.Valid() {
			break
		}

		// Битая база?
		if err := iterator.Error(); err != nil {
			chatList = []string{}

			return fmt.Errorf("Unable to get chat list: %w", err)
		}

		chatName := iterator.Key()

		if string(chatName) == "" {
			return errors.New("Got empty chat name when query " + chatListDBName)
		}

		chatList = append(chatList, string(chatName))

		iterator.Next()
	}

	return err
}

// AppendChatListDB добавляет чат в список известных боту. Можно добавлять чаты не глядя, дубликатов не добавится.
// Bot API нам такую информацию не даёт.
func AppendChatListDB(chatID string) error {
	var err error

	if slices.Contains(chatList, chatID) {
		return err
	}

	if err = StoreKV(chatListDB, chatID, "{}"); err != nil {
		return fmt.Errorf("Unable to append %s to %s: %w", chatID, chatListDBName, err)
	}

	// TODO: fill in defaults!

	// Не забываем актуализировать список чатов.
	chatList = append(chatList, chatID)
	slices.Sort(chatList)

	return err
}

// GetSetting достаёт настройку из БД с настройками.
func GetSetting(chatID string, setting string) string {
	var err error

	chatHash := sha256.Sum256([]byte(chatID))
	database := fmt.Sprintf("settings_db/%x", chatHash)

	// Если БД не открыта, откроем её.
	if _, ok := settingsDB[database]; !ok {
		var options pebble.Options
		// По дефолту ограничение ставится на мегабайты данных, а не на количество файлов, поэтому с дефолтными
		// настройками порождается огромное количество файлов. Умолчальное ограничение на количество файлов - 500 штук,
		// что нас не устраивает, поэтому немного снизим эту цифру до более приемлемых значений.
		options.L0CompactionFileThreshold = 8
		settingsDB[database], err = pebble.Open(Config.DataDir+"/"+database, &options)

		if err != nil {
			log.Errorf("Unable to open settings db, %s: %s\n", database, err)

			return ""
		}
	}

	value, err := FetchV(settingsDB[database], setting)

	// Если из базы ничего не вынулось, по каким-то причинам, то просто вернём пустую строку.
	if err != nil {
		switch {
		case errors.Is(err, pebble.ErrNotFound):
			log.Debugf("Unable to get value for %s: no record found in db %s", setting, database)
		case errors.Is(err, fs.ErrNotExist):
			log.Debugf("Unable to get value for %s: db dir %s does not exist", setting, database)
		case errors.Is(err, oserror.ErrNotExist):
			log.Debugf("Unable to get value for %s: db dir %s does not exist", setting, database)
		default:
			log.Errorf("Unable to get value for %s in db dir %s: %s", setting, database, err)
		}

		return ""
	}

	return value
}

// SaveSetting сохраняет настройку в БД с настройками.
func SaveSetting(chatID string, setting string, value string) error {
	var (
		chatHash = sha256.Sum256([]byte(chatID))
		database = fmt.Sprintf("settings_db/%x", chatHash)
		err      error
	)

	// Если БД не открыта, откроем её.
	if _, ok := settingsDB[database]; !ok {
		var options pebble.Options
		// По дефолту ограничение ставится на мегабайты данных, а не на количество файлов, поэтому с дефолтными
		// настройками порождается огромное количество файлов. Умолчальное ограничение на количество файлов - 500 штук,
		// что нас не устраивает, поэтому немного снизим эту цифру до более приемлемых значений.
		options.L0CompactionFileThreshold = 8

		settingsDB[database], err = pebble.Open(Config.DataDir+"/"+database, &options)

		if err != nil {
			log.Errorf("Unable to open settings db, %s: %s\n", database, err)

			return err
		}
	}

	if err := StoreKV(settingsDB[database], setting, value); err != nil {
		log.Errorf("Unable to save setting %s in database %s: %s", setting, database, err)
	}

	return err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
