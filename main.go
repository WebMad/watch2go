package main

import (
	"fmt"

	"net/http"

	"github.com/gorilla/websocket"

	"encoding/json"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[string]*websocket.Conn)

func main() {
	var settings Settings
	settings.load()

	var videoData VideoData
	videoData.load()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			return
		}

		if clients[conn.RemoteAddr().String()] == nil {
			clients[conn.RemoteAddr().String()] = conn
			conn.SetCloseHandler(func(code int, text string) error {
				if clients[conn.RemoteAddr().String()] != nil {
					delete(clients, conn.RemoteAddr().String())
				}
				return nil
			})
		}

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
				send(conn, msgType, videoData.asJSON())
			case "pause":
				videoData.IsPaused = true
				videoData.Time = packet["time"].(float64)
				fmt.Println("[info] пауза")
				broadcast(clients, msgType, videoData.asJSON())
			case "play":
				videoData.IsPaused = false
				videoData.Time = packet["time"].(float64)
				fmt.Println("[info] воспроизведение")
				broadcast(clients, msgType, videoData.asJSON())
			case "changeTime":
				fmt.Println("[info] время изменено")
				videoData.Time = packet["time"].(float64)
				broadcast(clients, msgType, videoData.asJSON())
			case "src":
				fmt.Println("[info] изменен источник")
				videoData.Src = packet["src"].(string)
				videoData.Time = 0
				broadcast(clients, msgType, videoData.asJSON())
			case "msg":
				fmt.Println("[info] новое сообщение")
				var msgPacket MsgPacket
				msgPacket.PacketType = "msg"
				msgPacket.Msg = packet["msg"].(string)
				msgPacket.Type = packet["msgType"].(string)
				broadcast(clients, msgType, msgPacket.asJSON())
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

	http.ListenAndServe(settings.IP+":"+settings.Port, nil)
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
