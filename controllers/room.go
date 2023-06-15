package controllers

import (
	"fmt"
	"sync"
)

var GlobalRoomManager *RoomManager

type Room struct {
	ID      int
	clients []*Client
	state   int // 表示状态,0表示房间未满人，1表示房间满人
}

type RoomManager struct {
	rooms  map[int]*Room // 房间ID到房间的映射
	mutex  sync.Mutex    // 互斥锁，保证并发安全
	nextID int           // 用于生成下一个房间的ID
	queue  []*Client     // 用于匹配的玩家队列
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[int]*Room),
		queue: make([]*Client, 0),
	}
}

// AddToQueue adds a client to the matchmaking queue.
func (rm *RoomManager) AddToQueue(client *Client) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	fmt.Println("Start AddToQueue")
	rm.queue = append(rm.queue, client)
	fmt.Println("Added to queue: ", client.UserInfo.Username)

	// If we have enough clients in the queue to start a new game...
	if len(rm.queue) >= 3 {
		clients := rm.queue[:3]
		rm.queue = rm.queue[3:]

		room := &Room{
			ID:      rm.nextID,
			clients: clients,
			state:   0,
		}

		// Send game start message to all clients in the room
		room.sendGameStartMessage()

		rm.rooms[rm.nextID] = room
		rm.nextID++

		fmt.Println("Created new room with ID: ", room.ID)
	}
}

func (rm *RoomManager) GetRoom(id int) (*Room, bool) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	room, ok := rm.rooms[id]
	return room, ok
}

func (room *Room) sendGameStartMessage() {
	// 在这里，你可能需要创建你的游戏开始消息，可能是一个结构体或者其他数据类型
	// 这里假设你已经有了一个创建游戏开始消息的函数 createGameStartMessage()
	// 该函数返回一个字节数组

	// 然后，向房间中的所有用户发送游戏开始的消息
	for _, client := range room.clients {
		sendResponse(nil, 14, client)
	}
}
