package service

import (
	"fmt"
	"sync"
)

// 定义一个游戏房间的结构体
type Room struct {
	ID      int    // 房间ID
	Players []int  // 房间内的玩家ID
	Tables  map[TableId]*Table
	state int  //表示状态,0表示房间未满人，1表示房间满人
}

// 定义一个房间管理器的结构体
type RoomManager struct {
	Rooms map[int]*Room // 房间ID到房间的映射
	Mutex sync.Mutex    // 互斥锁，保证并发安全
}

// 创建一个新的房间管理器
func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms: make(map[int]*Room),
	}
}

// 创建一个新的房间，并添加到房间管理器中
func (rm *RoomManager) CreateRoom(id int) *Room {
	rm.Mutex.Lock()         // 加锁
	defer rm.Mutex.Unlock() // 解锁

	// 检查房间ID是否已经存在
	if _, ok := rm.Rooms[id]; ok {
		fmt.Println("Room ID already exists:", id)
		return nil
	}

	// 创建一个新的房间
	room := &Room{
		ID:      id,
		Players: make([]int, 0),
	}

	// 添加到房间管理器中
	rm.Rooms[id] = room

	return room
}

// 根据房间ID获取房间
func (rm *RoomManager) GetRoom(id int) *Room {
	rm.Mutex.Lock()         // 加锁
	defer rm.Mutex.Unlock() // 解锁

	return rm.Rooms[id]
}

// 删除一个房间
func (rm *RoomManager) DeleteRoom(id int) {
	rm.Mutex.Lock()         // 加锁
	defer rm.Mutex.Unlock() // 解锁

	delete(rm.Rooms, id)
}

// 向一个房间添加一个玩家
func (rm *RoomManager) AddPlayer(playerID int) {
	rm.Mutex.Lock()         // 加锁
	defer rm.Mutex.Unlock() // 解锁

	// if room, ok := rm.Rooms[roomID]; ok {
	// 	room.Players = append(room.Players, playerID)
	// }
	//遍历所有房间，当有房间的长度小于3的时候，把这个玩家加入房间
	for _, room := range rm.Rooms {
		fmt.Printf("Room ID: %d, Players: %v\n", room.ID, room.Players)
		if len(room.Players) <3 { 
			room.Players = append(room.Players, playerID)
			break;
		}
	}
}

// 从一个房间移除一个玩家
func (rm *RoomManager) RemovePlayer(roomID, playerID int) {
	rm.Mutex.Lock()         // 加锁
	defer rm.Mutex.Unlock() // 解锁

	if room, ok := rm.Rooms[roomID]; ok {
		for i, p := range room.Players {
			if p == playerID {
				room.Players = append(room.Players[:i], room.Players[i+1:]...)
				break
			}
		}
	}
}

// 打印所有的房间信息，用于测试
func (rm *RoomManager) PrintRooms() {
	rm.Mutex.Lock()         // 加锁
	defer rm.Mutex.Unlock() // 解锁

	for _, room := range rm.Rooms {
		fmt.Printf("Room ID: %d, Players: %v\n", room.ID, room.Players)
	}
}
