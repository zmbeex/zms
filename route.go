package zms

import (
	"github.com/zmbeex/gkit"
	"reflect"
	"strings"
)

type Router struct {
	Code   string      `title:"服务编码"`
	Desc   string      `title:"服务描述"`
	Role   []string    `title:"访问角色"`
	Params interface{} `title:"参数模型"`
	Func   func(z *Zms)
}

type ParamsModel struct {
	Title        string
	Check        string
	Name         string
	DefaultValue string
	Type         string
}

func (r Router) Register() {
	Cache.ServerRouterMap[r.Code] = &r
	info := new(ServerInfo)
	info.Code = r.Code
	info.Desc = r.Desc
	info.Role = r.Role
	info.Params = make(map[string]*ParamsModel)
	if r.Params != nil {
		json, _ := gkit.GetJSONStruct(r.Params)
		json.Range(func(typeField reflect.StructField, value reflect.Value) {
			model := new(ParamsModel)
			model.Name = gkit.StringFirstLower(typeField.Name)
			model.Check = typeField.Tag.Get("check")
			model.DefaultValue = typeField.Tag.Get("defaultValue")
			model.Title = typeField.Tag.Get("title")
			model.Type = typeField.Type.Name()
			info.Params[model.Name] = model
		})
	}
	Cache.ServerInfo[r.Code] = info
	gkit.Info("注册路由：" + info.Code + " -> " + info.Desc + " -> " + strings.Join(info.Role, ","))
	gkit.Info(gkit.SetJson(info.Params))
}

type ServerInfo struct {
	Code   string                  `title:"服务编码"`
	Desc   string                  `title:"服务描述"`
	Role   []string                `title:"访问角色"`
	Params map[string]*ParamsModel `title:"参数模型"`
}

type Zurl struct {
	val string
}

func (z *Zurl) Add(val string) {
	z.val += val
}

func (z *Zurl) Push(key string, value interface{}) {
	z.val += "&" + key + "=" + gkit.ToString(value)
}

func (z *Zurl) String() string {
	return z.val
}

/// 运行
func Run() {
	c := new(Client)

	// 重新连接
	c.reconnect()
	// 处理消息
	go c.handleMessage()
	// 发送心跳包
	c.Hearbeat()
}
