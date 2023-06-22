package controllers

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/astaxie/beego/logs"
)

//var Clienter *Client

const (
	writeWait      = 1 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512

	RoleFarmer   = 0
	RoleLandlord = 1
)

type UserId int
type UserInfo struct {
	UserId   UserId
	Username string
	Role     int
}

// 定义一个client的结构体
type Client struct {
	UserInfo *UserInfo
	conn     net.Conn // TCP连接
	addr     string   // 客户端地址
	port     string   // 客户端端口
	Room     *Room
	Table    *Table
	SeatNum  int    //座位顺序
	Hand     []Card // 玩家的手牌
	LordPoint int //地主分，初始第一个玩家分为2，第二个分为1，第三个分为0
	IsLord int //是否是地主
}

// 定义一个接收器，实现和客户端建立TCP连接的功能
func ConnectClient() {
	// 拼接客户端地址和端口
	//client := fmt.Sprintf("%s:%d", c.addr, c.port)
	// 使用net.Listen函数监听TCP的地址和端口信息

	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		logs.Error("listen Error:%v", err)
		return // 监听失败，返回错误
	}
	fmt.Println("Server is listening...")

	defer ln.Close() // 延迟关闭监听

	// 使用ln.Accept函数等待客户端连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Error("accept Error%v", err)
			return // 连接失败，返回错误
		}
		fmt.Println("New client connected.")
		//c.conn = conn // 连接成功，将conn赋值给c.conn
		client := &Client{
			conn:     conn,
			UserInfo: new(UserInfo),
		}
		go HandleRequest(client)
	}
}

func (c *Client) sendMsg(msg []interface{}) {
	// 将msg转换为字节切片
	msgByte, err := json.Marshal(msg)
	if err != nil {
		logs.Error("send msg [%v] marsha1 err:%v", string(msgByte), err)
		return
	}

	//设置写超时时间
	err = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		logs.Error("send msg SetWriteDeadline [%v] err:%v", string(msgByte), err)
		return
	}
	//向TCP连接的套接字中写入数据
	_, err = c.conn.Write(msgByte)
	if err != nil {
		logs.Error("send msg [%v] write err:%v", string(msgByte), err)
		return
	}
}
