package zms

import (
	"github.com/zmbeex/gkit"
	"strings"
)

type ZmsRouter struct {
	Code   string      `title:"服务编码"`
	Desc   string      `title:"服务描述"`
	Role   []string    `title:"访问角色"`
	Params interface{} `title:"参数模型"`
	Func   func(z *zms.Zms)
}

func (r ZmsRouter) Register() {
	Cache.ServerRouterMap[r.Code] = &r
	info := new(ServerInfo)
	info.Code = r.Code
	info.Desc = r.Desc
	info.Role = r.Role
	info.Params = map[string]interface{}{
		"Params": r.Params,
	}
	Cache.ServerInfo[r.Code] = info
	gkit.Info("注册路由：" + info.Code + " -> " + info.Desc + " -> " + strings.Join(info.Role, ","))
}

type ServerInfo struct {
	Code   string                 `title:"服务编码"`
	Desc   string                 `title:"服务描述"`
	Role   []string               `title:"访问角色"`
	Params map[string]interface{} `title:"参数模型"`
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

func Run() {
	c := new(ZmsClient)

	// 重新连接
	c.reconnect()
	// 处理消息
	go c.handleMessage()
	// 发送心跳包
	c.Hearbeat()
}
