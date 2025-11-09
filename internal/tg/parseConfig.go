package tg

import (
	"encoding/json"
	"fmt"
	"os"

	"aleesa-telegram-go/internal/log"

	"github.com/hjson/hjson-go"
)

// parseConfig разбирает и валидирует даденный конфиг.
func parseConfig(path string) (MyConfig, error) {
	fileInfo, err := os.Stat(path)

	// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём.
	if err != nil {
		return MyConfig{}, err
	}

	// Конфиг-файл длинноват для конфига, попробуем следующего кандидата.
	if fileInfo.Size() > 65535 {
		err := fmt.Errorf("config file %s is too long for config, skipping", path)

		return MyConfig{}, err
	}

	buf, err := os.ReadFile(path)

	// Не удалось прочитать.
	if err != nil {
		return MyConfig{}, err
	}

	// Исходя из документации, hjson какбы умеет парсить "кривой" json, но парсит его в map-ку.
	// Интереснее на выходе получить структурку: то есть мы вначале конфиг преобразуем в map-ку, затем эту map-ку
	// сериализуем в json, а потом json превращааем в структурку. Не очень эффективно, но он и не часто требуется.
	var (
		sampleConfig MyConfig
		tmp          map[string]any
	)

	err = hjson.Unmarshal(buf, &tmp)

	// Не удалось распарсить.
	if err != nil {
		return MyConfig{}, err
	}

	tmpjson, err := json.Marshal(tmp)

	// Не удалось преобразовать map-ку в json.
	if err != nil {
		return MyConfig{}, err
	}

	if err := json.Unmarshal(tmpjson, &sampleConfig); err != nil {
		return MyConfig{}, err
	}

	// Валидируем значения из конфига.
	if sampleConfig.Redis.Server == "" {
		sampleConfig.Redis.Server = "localhost"

		log.Infof("Redis server is not defined in config, using localhost")
	}

	if sampleConfig.Redis.Port == 0 {
		sampleConfig.Redis.Port = 6379
	}

	if sampleConfig.Redis.Channel == "" {
		err := fmt.Errorf("channel field in config file %s must be set", path)

		return MyConfig{}, err
	}

	if sampleConfig.Redis.MyChannel == "" {
		err := fmt.Errorf("my_channel field in config file %s must be set", path)

		return MyConfig{}, err
	}

	if sampleConfig.Telegram.Token == "" {
		return MyConfig{}, fmt.Errorf("telegram.token field in config file %s must be set", path)
	}

	if sampleConfig.Loglevel == "" {
		sampleConfig.Loglevel = "info"
	}

	// sampleConfig.Log = "" if not set..

	if sampleConfig.Csign == "" {
		err := fmt.Errorf("csign field in config file %s must be set", path)

		return MyConfig{}, err
	}

	if sampleConfig.ForwardsMax == 0 {
		sampleConfig.ForwardsMax = forwardMax
	}

	if sampleConfig.DataDir == "" {
		return MyConfig{}, fmt.Errorf("data_dir field in config file %s must be set", path)
	}

	return sampleConfig, err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
