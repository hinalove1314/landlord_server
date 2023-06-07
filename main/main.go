package main

import (
	"my_landlord/service"

	"github.com/astaxie/beego/logs"
)

func main() {
	//c := new(service.Client)

	err := initConf()
	if err != nil {
		logs.Error("init conf err:%v", err)
		return
	}
	defer func() {
		if gameConf.Db != nil {
			err = gameConf.Db.Close()
			if err != nil {
				logs.Error("main close sqllite db err :%v", err)
			}
		}
	}()
	err = initSec()
	if err != nil {
		logs.Error("init sec err:%v", err)
		return
	}

	//client:=&service.Client{}
	//无限循环来建立TCP连接?
	service.ConnectClient() //建立TCP连接

	// for {
	// 	conn, err := listen.Accept()
	// 	if err != nil {
	// 		fmt.Println("Error accepting connection:", err.Error())
	// 		continue
	// 	}

	// 	go handleRequest(conn)

	//controllers.StartConnect()
}
