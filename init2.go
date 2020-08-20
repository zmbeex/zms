package zms

import "github.com/satori/go.uuid"

var Cache struct {
	ServerRouterMap map[string]*ZmsRouter
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

	Cache.ServerRouterMap = make(map[string]*ZmsRouter)
	Cache.Uuid = uuid.NewV4().String()
	Cache.ServerInfo = make(map[string]*ServerInfo)
}
