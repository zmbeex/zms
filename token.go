package zms

import (
	"github.com/zmbeex/dao/tredis"
	"github.com/zmbeex/gkit"
	"time"
)

var tokenPrefix = "zms.token_"

type Token struct {
	Id       int64  `title:"用户id"`
	Platform int    `title:"平台"`    // 平台 1androi 2ios 3web
	Time     int64  `title:"授权有效期"` // 时间戳
	Code     string `title:"授权编码"`
	Device   string `title:"设备码"`
	err      error
	access   string
}

// 创建token
func SetToken(userId int64, platform int, code string, device string) *Token {
	t := new(Token)
	if userId == 0 {
		return t
	}
	t.Id = userId
	t.Platform = platform
	t.Time = time.Now().Unix()
	t.Code = code
	t.Device = device
	return t
}

// 每次生成一个新的access,并以此为键，存储data的值
func (t *Token) SetAccess() string {
	s := gkit.SetJson(map[string]interface{}{
		"Id":       t.Id,
		"Platform": t.Platform,
		"Time":     time.Now().UnixNano(),
		"Code":     t.Code,
		"Device":   t.Device,
	})
	t.access = gkit.GetSHA(s)

	err := tredis.SetRedis(
		tokenPrefix+t.access,
		gkit.ToString(t.Id)+"#"+s,
		time.Duration(Cache.Set.AccessTime)*time.Second,
	)
	if err != nil {
		gkit.Warn(err.Error())
		return ""
	}

	err = tredis.SetRedis(
		gkit.ToString(t.Id)+"#"+s,
		tokenPrefix+t.access,
		time.Duration(Cache.Set.AccessTime)*time.Second,
	)
	if err != nil {
		gkit.Warn(err.Error())
		return ""
	}
	return t.access
}

// 删除access
func (t *Token) DelAccess(access string) {
	flag := tredis.GetRedis(tokenPrefix + access)
	_ = tredis.SetRedis(access, flag, 30*time.Second)
	_ = tredis.SetRedis(flag, access, 30*time.Second)
}

func (t *Token) SetToken() string {
	s := gkit.SetJson(t)
	s = gkit.GetSHA(s)

	err := tredis.SetRedis(tokenPrefix+s, gkit.SetJson(t), time.Duration(Cache.Set.TokenTime)*time.Second)
	if err != nil {
		gkit.Warn(err.Error())
		return ""
	}
	return s
}

func GetToken(token string) *Token {
	t := new(Token)
	s := tredis.GetRedis(tokenPrefix + token)
	if s == "" {
		return nil
	}
	err := gkit.GetJson(s, t)
	if err != nil {
		return nil
	}
	return t
}

func GetAccess(access string) string {
	s := tredis.GetRedis(tokenPrefix + access)
	return s
}