package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hjson/hjson-go"
	log "github.com/sirupsen/logrus"
)

// parseConfig разбирает и валидирует даденный конфиг.
func parseConfig(path string) (myConfig, error) {
	fileInfo, err := os.Stat(path)

	// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём.
	if err != nil {
		return myConfig{}, err
	}

	// Конфиг-файл длинноват для конфига, попробуем следующего кандидата.
	if fileInfo.Size() > 65535 {
		err := fmt.Errorf("Config file %s is too long for config, skipping", path)

		return myConfig{}, err
	}

	buf, err := os.ReadFile(path)

	// Не удалось прочитать.
	if err != nil {
		return myConfig{}, err
	}

	// Исходя из документации, hjson какбы умеет парсить "кривой" json, но парсит его в map-ку.
	// Интереснее на выходе получить структурку: то есть мы вначале конфиг преобразуем в map-ку, затем эту map-ку
	// сериализуем в json, а потом json преврщааем в стркутурку. Не очень эффективно, но он и не часто требуется.
	var (
		sampleConfig myConfig
		tmp          map[string]interface{}
	)

	err = hjson.Unmarshal(buf, &tmp)

	// Не удалось распарсить.
	if err != nil {
		return myConfig{}, err
	}

	tmpjson, err := json.Marshal(tmp)

	// Не удалось преобразовать map-ку в json.
	if err != nil {
		return myConfig{}, err
	}

	if err := json.Unmarshal(tmpjson, &sampleConfig); err != nil {
		return myConfig{}, err
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
		err := fmt.Errorf("Channel field in config file %s must be set", path)

		return myConfig{}, err
	}

	if sampleConfig.Redis.MyChannel == "" {
		err := fmt.Errorf("My_channel field in config file %s must be set", path)

		return myConfig{}, err
	}

	if sampleConfig.Telegram.Token == "" {
		return myConfig{}, fmt.Errorf("telegram.token field in config file %s must be set", path)
	}

	if sampleConfig.Loglevel == "" {
		sampleConfig.Loglevel = "info"
	}

	// sampleConfig.Log = "" if not set

	if sampleConfig.Csign == "" {
		err := fmt.Errorf("Csign field in config file %s must be set", path)

		return myConfig{}, err
	}

	if sampleConfig.ForwardsMax == 0 {
		sampleConfig.ForwardsMax = forwardMax
	}

	if sampleConfig.DataDir == "" {
		return myConfig{}, fmt.Errorf("Data_dir field in config file %s must be set", path)
	}

	return sampleConfig, err
}
