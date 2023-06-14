package controllers

import (
	"sync"
)

// 这样声明结构体中的变量是为了让这个变量在room.go中能够找得到
type TableId int

const (
	GameWaitting = iota
	GameCallScore
	GamePlaying
	GameEnd
)

type Table struct {
	Lock      sync.RWMutex
	TableID   TableId
	State     bool
	PlayerNum int
}

func (table *Table) joinTable() {
	table.Lock.Lock()
	defer table.Lock.Unlock()

	//当桌子的人数大于2的时候，报错

	//当加入桌子的是同一个人时，报错

}

// 同步用户信息
func (table *Table) syncUser() {

}
