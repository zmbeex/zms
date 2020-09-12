package zms

import (
	uuid "github.com/satori/go.uuid"
	"github.com/zmbeex/gkit"
	"time"
)

// 参数
type Params struct {
	Code   string `title:"服务编码"`
	UserId int64  `title:"用户id"`
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

const (
	ROLE_PUBLIC = "zms_public"
	ROLE_ADMIN  = "zms_admin"
)

var Cache struct {
	ServerRouterMap map[string]*Router
	ServerInfo      map[string]*ServerInfo
	Uuid            string
	Set             *Setting
}

type Setting struct {
	UserName      string `title:"账号" defaultValue:"dev"`
	Password      string `title:"密码" defaultValue:"qwertyuiop3466f"`
	GatewayHost   string `title:"网关" defaultValue:"ws://localhost:8088"`
	CallHost      string `title:"网关" defaultValue:"http://localhost:8088"`
	ServerInfoKey string `title:"服务信息加密密钥" defaultValue:"1234567812345678"`
	AccessTime    int64  `title:"access 有效时间，默认10分钟" defaultValue:"600"`
	TokenTime     int64  `title:"access 有效时间, 默认10天" defaultValue:"864000"`
}

func init() {
	set := new(Setting)
	gkit.InitSetting("zms", set, "zms客户端", func() {
		Cache.Set = set
		Cache.ServerRouterMap = make(map[string]*Router)
		Cache.Uuid = uuid.NewV4().String()
		Cache.ServerInfo = make(map[string]*ServerInfo)
	})
}

// 调用其他服务
func CallServer(code string, params string, userId int64) *Result {
	p := new(Result)
	p.Code = code
	func() {
		defer func() {
			r := recover()
			if r != nil {
				p.Note = "调用服务异常"
				p.Status = -888
				p.Note = gkit.ToString(r)
			}
		}()
		timeStamp := time.Now().Unix()
		sign := "zms.call|" + gkit.ToString(timeStamp) + "|" + code + "|" + params
		sign = gkit.GetSHA(sign)
		data := gkit.HttpPOST(Cache.Set.CallHost+"/call", map[string]string{
			"code":      code,
			"sign":      sign,
			"params":    params,
			"userId":    gkit.ToString(userId),
			"clientUid": Cache.Set.UserName + "#" + Cache.Uuid,
			"timeStamp": gkit.ToString(time.Now().Unix()),
		}).ToString()

		err := gkit.GetJson(data, p)
		if err != nil {
			gkit.Error("CallServer执行失败")
			gkit.Error(err.Error())
		}
	}()
	return p
}
