package services

import (
	"context"
	"encoding/json"

	"fmt"
	"piga-go/models"
	"time"

	"github.com/astaxie/beego"
)

func init() {

	// runList = make(map[string]*models.Task)
	go sendBroadCastEvent()
}

// var runList map[string]*models.Task

func sendBroadCastEvent() {

	for {
		select {
		case event := <-models.BroadCast:
			beego.Debug(fmt.Sprintf("获取新事件 %#v", event))
			body, err := json.Marshal(event.RequestEvent)
			if err != nil {
				event.Response <- models.JsonFn(models.Error, "error", fmt.Sprintf("任务json解析错误 原因: %#v", err))
				break
			}

			task := models.Task{
				Uuid:     event.RequestEvent.Uuid,
				Response: event.Response,
			}
			if event.RequestEvent.Timeout > 15 || event.RequestEvent.Timeout < 1 {
				event.RequestEvent.Timeout = 5
			}
			ctx := context.Background()
			task.Ext, task.Cancel = context.WithTimeout(ctx, time.Duration(event.RequestEvent.Timeout)*time.Second)
			go task.Run(event.WsConn, body)
		case <-time.Tick(time.Second * 5):
			beego.Debug("每隔一分钟进行一次打印", time.Now(), cap(models.BroadCast), len(models.BroadCast))
		}
	}
	fmt.Println("循环结束...")

}
