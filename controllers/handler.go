package controllers

import (
	"encoding/binary"
	"fmt"
	"io"
	"my_landlord/common"

	"encoding/json"

	"github.com/astaxie/beego/logs"
)

//var rm RoomManager

//var matchQueue = &MatchQueue{}

var netCode struct {
	Value int `json:"value"`
}

func HandleRequest(client *Client) {
	fmt.Println("HandleRequest")

	for {
		// 先读取长度
		lengthBuf := make([]byte, 4)
		_, err := io.ReadFull(client.conn, lengthBuf)
		if err != nil {
			logs.Error("read data head Error")
			return
		}
		length := binary.LittleEndian.Uint32(lengthBuf)
		fmt.Println("length=", length)

		// 读取netCode
		netCodeBuf := make([]byte, 4)
		_, err = io.ReadFull(client.conn, netCodeBuf)
		if err != nil {
			logs.Error("read data head Error")
			return
		}
		netCode := binary.LittleEndian.Uint32(netCodeBuf)
		fmt.Println("netCode:", netCode)

		// 读取消息，注意长度需要减去已经读取的netCode的长度
		msgBuf := make([]byte, length-4)
		_, err = io.ReadFull(client.conn, msgBuf)
		if err != nil {
			logs.Error("read data body Error")
			return
		}
		fmt.Println("debug_msgbuf")
		fmt.Println("msgBuf=", string(msgBuf))

		// 现在msgBuf就是完整的json数据，可以直接进行解析
		// err = json.Unmarshal(msg, &data)

		switch netCode {
		case common.ReqCreat:
			Register(msgBuf, client)
		case common.ReqLogin:

		case common.ReqRoomList:
			fmt.Println("common.ReqRoomList")
			GlobalRoomManager.AddToQueue(client)
		case common.ReqJoinRoom:
			// rm.PrintRooms()
			// rm.Mutex.Lock()
			// defer rm.Mutex.Unlock()
			// rm.AddPlayer(int(client.UserInfo.UserId))
			// rm.PrintRooms()
		}
	}
}

func sendResponse(response interface{}, netCode int, client *Client) {
	fmt.Println("sendResponse start")

	jsonData, err := json.Marshal(response)
	if err != nil {
		logs.Error("Error marshalling JSON:", err)
		return
	}

	// prepare the data to send
	length := len(jsonData)
	fmt.Printf("length: %v\n", length)

	data := make([]byte, 4+4+len(jsonData))
	fmt.Printf("data: %v\n", data)

	binary.LittleEndian.PutUint32(data[0:4], uint32(length))
	binary.LittleEndian.PutUint32(data[4:8], uint32(netCode))
	copy(data[8:], jsonData)

	fmt.Printf("data=%v", data)
	_, err = client.conn.Write(data)
	if err != nil {
		logs.Error("Error sending data:", err)
	} else {
		fmt.Println("Data sent successfully")
	}
}
