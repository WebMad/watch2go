package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// VideoData структура для хранения состаяния воспроизведения видео
type VideoData struct {
	PacketType string  `json:"packetType"`
	Src        string  `json:"src"`
	Time       float64 `json:"time"`
	IsPaused   bool    `json:"isPaused"`
}

func (videoData *VideoData) load() {
	dataFile, err := os.Open("data.json")

	defer func() {
		err := dataFile.Close()
		if err != nil {
			fmt.Println("[error] ошибка закрытия файла")
		}
	}()

	if err != nil {
		fmt.Println("[error] ошибка открытия файла сохранения")
		return
	}

	byteValue, _ := ioutil.ReadAll(dataFile)
	json.Unmarshal([]byte(byteValue), &videoData)

	videoData.IsPaused = true
	videoData.PacketType = "videoData"
}

func (videoData *VideoData) asJSON() string {
	encodedData, _ := json.Marshal(videoData)
	return string(encodedData)
}
