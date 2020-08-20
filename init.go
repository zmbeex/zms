package zms

import uuid "github.com/satori/go.uuid"

// 参数
type Params struct {
	Code   string `title:"服务编码"`
	Token  string `title:"认证签名"`
	Params string `title:"请求入参"`
	Uuid   string `title:"唯一标识"`
}

// 返回数据
type Result struct {
	Uuid   string                 `title:"唯一标识"`
	Code   string                 `title:"服务编码"`
	Status int                    `title:"状态"`
	Note   string                 `title:"提示信息"`
	Data   map[string]interface{} `title:"返回结果/任意类型"`
	List   []interface{}          `title:"返回结果，数组"`
}

var Cache struct {
	ServerRouterMap map[string]*Router
	ServerInfo      map[string]*ServerInfo
	Uuid            string
	Set             *Setting
}

type Setting struct {
	UserName      string `title:"账号"`
	Password      string `title:"密码"`
	GatewayHost   string `title:"网关"`
	ServerInfoKey string `title:"服务信息加密密钥"`
}

func init() {
	set := new(Setting)
	set.UserName = "dev"
	set.Password = "qwertyuiop3466f"
	set.GatewayHost = "wss://zmbeex.com/gateway/nest"
	//set.GatewayHost = "ws://192.168.0.108:8088/nest"
	set.ServerInfoKey = "1234567812345678"
	Cache.Set = set

	Cache.ServerRouterMap = make(map[string]*Router)
	Cache.Uuid = uuid.NewV4().String()
	Cache.ServerInfo = make(map[string]*ServerInfo)
}
