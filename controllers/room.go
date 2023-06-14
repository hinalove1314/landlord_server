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
