package common

import (
	"database/sql"
)

var GameConfInfo GameConf
var LoginResponseInfo LoginResponse

type GameConf struct {
	HttpPort int
	LogPath  string
	DbPath   string

	LogLevel string
	Db       *sql.DB
}

type LoginResponse struct {
	IsRegisted int    `json:"m_isRegisted"`
	IsLogin    int    `json:"m_isLogin"`
	Account    string `json:"m_dataAccount"`
	Password   string `json:"m_dataPassword"`
}
