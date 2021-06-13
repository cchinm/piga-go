package models

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var mux sync.RWMutex

type TaskFunc func(ws *websocket.Conn)

type Task struct {
	Uuid     string
	Response chan string
	Ext      context.Context
	Cancel   context.CancelFunc
	fn       TaskFunc
}

func (t *Task) Run(ws *websocket.Conn, body []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover ", err)
			// t.Response <- fmt.Sprintf("读取消息过程中发送错误 原因 %#v", err)
			select {
			case t.Response <- JsonFn(Error, "error", fmt.Sprintf("%#v", err)):
				fmt.Println("recovery", err)
			case <-time.After(3 * time.Second):
				fmt.Println("recovery fail", err)
			}
		}
	}()

	if err := ws.WriteMessage(websocket.TextMessage, body); err != nil {
		switch err {
		case websocket.ErrCloseSent:
			AutoLeave(t.Uuid, "")
		default:
			AutoLeave(t.Uuid, "")
		}
		t.Response <- JsonFn(Error, "error", fmt.Sprintf("任务执行错误 原因: %#v", err))
		return
	}

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Response <- JsonFn(Error, "error", fmt.Sprintf("读取消息过程中发送错误 原因 %#v", err))
	} else {
		t.Response <- JsonFn(Success, "succ", string(p))
	}

}

func AutoLeave(uuid, group string) {
	mux.Lock()
	defer mux.Unlock()
	fmt.Println("删除分组:", group, "UUID: ", uuid)
	if _, ok := WsPool[uuid]; ok {
		delete(WsPool, uuid)
	}
	if _, ok := WsConnMap[uuid]; ok {
		delete(WsConnMap, uuid)
	}
	fmt.Println("清除成功......")

}

func JoinWs(ws *websocket.Conn, uuid, group string) {
	mux.RLock()
	defer mux.RUnlock()

	WsPool[uuid] = ws
	WsConnMap[uuid] = group
}
