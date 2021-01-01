package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Settings стуктура для настроек ip и порта сервера
type Settings struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func (settings *Settings) load() {
	settingsFile, err := os.Open("settings.json")
	if err != nil {
		fmt.Println("[error] ошибка загрузки настроек сервера")
		return
	}

	defer func() {
		if err := settingsFile.Close(); err != nil {
			fmt.Println("[error] ошибка закрытия файла настроек")
		}
	}()

	settingsByteValue, _ := ioutil.ReadAll(settingsFile)

	json.Unmarshal([]byte(settingsByteValue), &settings)
}
