package services

import (
	"fmt"
	"piga-go/models"
	"sync"
	"time"
)

var mux sync.RWMutex

func AppSearchByGroupName(group string) models.ReturnTpl {
	data := make(map[string]int)
	for key := range models.WsConnMap {
		if models.WsConnMap[key] == group {
			data[key] = 1
		}
	}
	return models.Success("succ", data)
}

func AppSearchAllGroup() models.ReturnTpl {
	data := make(map[string]int)
	for key := range models.WsConnMap {
		data[models.WsConnMap[key]] = data[models.WsConnMap[key]] + 1
	}
	return models.Success("succ", data)
}

func AppExecuteRemote(m *models.ExecuteEvent) models.ReturnTpl {

	if ws, ok := models.WsPool[m.Uuid]; !ok {
		return models.Error("链接不存在 请重试其他客户端", nil)
	} else {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		ret := make(chan string)
		fmt.Println("uuid", m.Uuid)
		broadEvent := models.BroadCastEvent{
			RequestEvent: m,
			WsConn:       ws,
			Response:     ret,
		}
		select {
		case models.BroadCast <- broadEvent:
			break
		case <-time.After(15 * time.Second):
			return models.Error("事件初始化超时", m)
		}
		select {

		case s := <-ret:
			result := models.Result{
				RequestEvent: m,
				Response:     s,
			}
			// close(ret)
			return models.Success("succ", result)
		case <-time.After(15 * time.Second):
			// close(ret)
			defer ws.Close()
			return models.Error("超时退出调用", m)
		}
	}
}
