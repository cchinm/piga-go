package controllers

import (
	"net/http"
	"piga-go/models"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type WebsocketController struct {
	beego.Controller
}

func (this *WebsocketController) GetTpl() {
	group := this.GetString("group")
	if group == "" {
		this.Data["json"] = "group不能为空"
		this.ServeJSON()
		return
	}
	this.Data["Uname"] = group
	this.TplName = "welcome.html"
	this.Render()
}

func (this *WebsocketController) Connect() {
	autoLeave := this.GetString("auto")
	uuid := this.GetString("uuid")
	group := this.GetString("group")
	var ws *websocket.Conn
	var err error
	if autoLeave == "true" {
		defer models.AutoLeave(uuid, group)
	}

	ws, err = websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 512, 512)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "这不是一个ws连接", 400)
		return
	} else if err != nil {
		beego.Error("ws建立失败:", err)
		return
	} else {
		models.JoinWs(ws, uuid, group)
		return
	}
}
