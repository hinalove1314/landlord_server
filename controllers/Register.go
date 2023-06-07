package controllers

import (
	"database/sql"
	"encoding/json"
	"my_landlord/common"
	"my_landlord/service"
	"time"

	"github.com/astaxie/beego/logs"
)

var (
	LogInfo = &common.LoginResponseInfo
)

// 可以直接传一个连接，也可以传一个buffer，待修改
func Register(msg []byte, client *service.Client) {
	err := json.Unmarshal(msg, &LogInfo)
	if err != nil {
		logs.Error("Error parsing JSON:", err)
		return
	}

	logs.Error("Account=%v", LogInfo.Account)
	logs.Error("Password=%v", LogInfo.Password)

	var username, password string
	err = common.GameConfInfo.Db.QueryRow("SELECT username FROM account where username=?", LogInfo.Account).Scan(&username, &password)

	// Use LoginResponseInfo instead of LoginResponse
	LogInfo.Account = LogInfo.Account
	LogInfo.Password = LogInfo.Password

	if err != nil && err != sql.ErrNoRows { //账号已存在
		print("账号已存在")
		LogInfo.IsRegisted = 0 //0表示注册失败
	} else {
		now := time.Now().Format("2006-01-02 15:04:05")
		_, err = common.GameConfInfo.Db.Exec("Insert INTO account (username,password,created_date,updated_date) values(?,?,?,?)", username, password, now, now)
		if err != nil {
			logs.Error("Insert account Error:%v", err)
			return
		}
		print("账号不存在，成功注册!")
		LogInfo.IsLogin = 1 //1表示登录成功
	}

	// 发送响应数据
	sendResponse(LogInfo, 2, client)
}
