package controllers

import (
	"encoding/json"
	"piga-go/models"
	"piga-go/services"

	"github.com/astaxie/beego"
)

type AppController struct {
	beego.Controller
}

// @Title Remote Get GroupUserList
// @Description 根据分组标签获取相应的用户链接
// @Param	group query true 分组名称
// @Success 200 {object} []models.ReturnTpl
// @Failure 403 空 无内容
// @router /groupUserList [get]
func (this *AppController) GroupUserList() {
	group := this.GetString("group")
	if group == "" {
		this.Data["json"] = models.Error("group参数为空", nil)
		this.ServeJSON()
		return
	}
	this.Data["json"] = services.AppSearchByGroupName(group)
	this.ServeJSON()
	return
}

// @Title Remote Get GroupList
// @Description 获取所有标签分组
// @Success 200 {object} []models.ReturnTpl
// @Failure 403 空 无内容
// @router /groupList [get]
func (this *AppController) GroupList() {
	this.Data["json"] = services.AppSearchAllGroup()
	this.ServeJSON()
	return
}

// @Title 远程调用函数
// @Description 根据id调用相应远程服务
// @Param	body 	body 	models.ExecuteEvent	true		"payload参数"
// @Success 200 {object} []models.ReturnTpl
// @Failure 403 空 无内容
// @router /execute [post]
func (this *AppController) Execute() {
	var event models.ExecuteEvent
	json.Unmarshal(this.Ctx.Input.RequestBody, &event)
	this.Data["json"] = services.AppExecuteRemote(&event)
	this.ServeJSON()
}
