package main

import "encoding/json"

// MsgPacket пакет отправляемый для вывода у пользователя
type MsgPacket struct {
	PacketType string `json:"packetType"`
	Type       string `json:"type"`
	Msg        string `json:"msg"`
}

func (msgPacket *MsgPacket) asJSON() string {
	encodedData, _ := json.Marshal(msgPacket)
	return string(encodedData)
}
