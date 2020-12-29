package main

import (
	"fmt"

	"net/http"

	"github.com/gorilla/websocket"

	"encoding/json"

	"os"

	"io/ioutil"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type VideoData struct {
	Src      string  `json:"src"`
	Time     float64 `json:"time"`
	IsPaused bool    `json:"isPaused"`
}

func (videoData *VideoData) asJson() string {
	encodedData, _ := json.Marshal(videoData)
	return string(encodedData)
}

var clients = make(map[string]*websocket.Conn)

func main() {
	dataFile, err := os.Open("data.json")
	var videoData VideoData

	if err != nil {
		fmt.Println("[error] ошибка открытия файла сохранения")
		return
	}

	byteValue, _ := ioutil.ReadAll(dataFile)
	json.Unmarshal([]byte(byteValue), &videoData)

	videoData.IsPaused = true

	defer func() {
		err := dataFile.Close()
		if err != nil {
			fmt.Println("[error] ошибка закрытия файла")
		}
	}()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		if clients[conn.RemoteAddr().String()] == nil {
			conn.SetCloseHandler(func(code int, text string) error {
				if clients[conn.RemoteAddr().String()] != nil {
					delete(clients, conn.RemoteAddr().String())
				}
				return nil
			})
			clients[conn.RemoteAddr().String()] = conn
		}
		fmt.Println(clients)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var packet map[string]interface{}

			err = json.Unmarshal([]byte(msg), &packet)
			if err != nil {
				fmt.Println("[error] ошибка декодирования json")
			}

			switch packet["type"] {
			case "getData":
				fmt.Println("[info] запрос на получение данных")
				send(conn, msgType, videoData.asJson())
			case "pause":
				videoData.IsPaused = true
				videoData.Time = packet["time"].(float64)
				fmt.Println("[info] пауза")
				broadcast(clients, msgType, videoData.asJson())
			case "play":
				videoData.IsPaused = false
				videoData.Time = packet["time"].(float64)
				fmt.Println("[info] воспроизведение")
				broadcast(clients, msgType, videoData.asJson())
			case "changeTime":
				fmt.Println("[info] время изменено")
				videoData.Time = packet["time"].(float64)
				broadcast(clients, msgType, videoData.asJson())
			case "src":
				fmt.Println("[info] изменен источник")
				videoData.Src = packet["src"].(string)
				videoData.Time = 0
				broadcast(clients, msgType, videoData.asJson())
			default:
				fmt.Println("[error] неизвестный тип сообщения:", packet["type"])
			}
			//dataFile.WriteString(videoData.asJson())
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("frontend/assets/"))))

	http.ListenAndServe(":25565", nil)
}

func send(conn *websocket.Conn, msgType int, msg string) {
	err := conn.WriteMessage(msgType, []byte(msg))
	if err != nil {
		fmt.Println(err)
		fmt.Println("[error] ошибка отправки сообщения")
		return
	}
}

func broadcast(clients map[string]*websocket.Conn, msgType int, msg string) {
	for _, conn := range clients {
		send(conn, msgType, msg)
	}
}
