package service

import (
	"encoding/binary"
	"fmt"
	"io"
	"my_landlord/common"
	"my_landlord/controllers"

	"github.com/astaxie/beego/logs"
)

var rm RoomManager

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
		msgBuf := make([]byte, length)
		_, err = io.ReadFull(client.conn, msgBuf)
		if err != nil {
			logs.Error("read data body Error")
			return
		}

		fmt.Println("msgBuf=", string(msgBuf))

		// 现在msgBuf就是完整的json数据，可以直接进行解析
		// err = json.Unmarshal(msg, &data)

		switch netCode {
		case common.ReqCreat:
			controllers.Register(msgBuf, err)
		case common.ReqLogin:

		case common.ReqRoomList:

		case common.ReqJoinRoom:
			rm.PrintRooms()
			rm.Mutex.Lock()
			defer rm.Mutex.Unlock()
			rm.AddPlayer(int(client.UserInfo.UserId))
			rm.PrintRooms()
		}
	}

}
