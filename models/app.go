package models

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type ReturnTpl struct {
	Code    int
	Message string
	Data    interface{}
}

const (
	SUCCESS = 10000*iota + 10000
	RISK
	UNLESS
	ERROR
	UNKNOWN
)

// uuid: ws.Conn
var WsPool map[string]*websocket.Conn

// uuid: group
var WsConnMap map[string]string
var BroadCast chan BroadCastEvent
var HealthBroadCast chan BroadCastEvent

func init() {
	WsConnMap = make(map[string]string)
	// WsConnMap["1"] = "a"
	// WsConnMap["2"] = "a"
	// WsConnMap["4"] = "b"
	// WsConnMap["3"] = "c"

	WsPool = make(map[string]*websocket.Conn)
	// WsPool["1"] = new(websocket.Conn)

	BroadCast = make(chan BroadCastEvent, 100)
	HealthBroadCast = make(chan BroadCastEvent, 100)
}

type ExecuteEvent struct {
	Param        string
	Method       string
	Timeout      int
	Uuid         string
	ForceExecute bool
}

type Result struct {
	RequestEvent *ExecuteEvent
	Response     string
}

type BroadCastEvent struct {
	RequestEvent *ExecuteEvent
	WsConn       *websocket.Conn `json:-`
	Response     chan string     `json:-`
}

func Success(message string, data interface{}) ReturnTpl {
	ret := ReturnTpl{
		Code:    SUCCESS,
		Message: message,
		Data:    data,
	}
	return ret
}

func Error(message string, data interface{}) ReturnTpl {
	ret := ReturnTpl{
		Code:    ERROR,
		Message: message,
		Data:    data,
	}
	return ret
}

type tplFunc func(message string, data interface{}) ReturnTpl

func JsonFn(fn tplFunc, message string, data interface{}) string {
	tpl := fn(message, data)
	p, _ := json.Marshal(&tpl)
	return string(p)
}
