package zms

import (
	"github.com/zmbeex/dao/tredis"
	"github.com/zmbeex/gkit"
	"strings"
	"time"
)

var tokenPrefix = "zms.token_"
var accessPrefix = "zms.access_"
var userRolePrefix = "zms.role_"
var userTokenInfoPrefix = "zms.user.info_"

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
func SetToken(userId int64, platform int) *Token {
	t := new(Token)
	if userId == 0 {
		return t
	}
	t.Id = userId
	t.Platform = platform
	t.Time = time.Now().Unix()
	return t
}

// 每次生成一个新的access,并以此为键，存储data的值
func (t *Token) SetAccess() (string, int64) {
	s := gkit.SetJson(map[string]interface{}{
		"Id":       t.Id,
		"Platform": t.Platform,
		"Code":     t.Code,
		"Device":   t.Device,
		"Time":     time.Now().UnixNano(),
	})
	t.access = gkit.GetSHA(s)

	err := tredis.SetRedis(
		accessPrefix+t.access,
		s,
		time.Duration(Cache.Set.AccessTime)*time.Second,
	)
	// 记录
	_ = SetTokenAllData(t.Id, accessPrefix+t.access)
	if err != nil {
		gkit.Warn(err.Error())
		return "", 0
	}
	return t.access, Cache.Set.AccessTime
}

// 删除access
func (t *Token) DelAccess(access string) {
	flag := tredis.GetRedis(tokenPrefix + access)
	_ = tredis.SetRedis(tokenPrefix+access, flag, 60*time.Second)
}

func (t *Token) SetToken() (string, int64) {
	s := gkit.SetJson(t)
	token := gkit.GetSHA(s)
	err := tredis.SetRedis(tokenPrefix+token, s, time.Duration(Cache.Set.TokenTime)*time.Second)
	// 记录
	_ = SetTokenAllData(t.Id, tokenPrefix+token)
	if err != nil {
		gkit.Warn(err.Error())
		return "", 0
	}
	return token, Cache.Set.TokenTime
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

func GetAccess(access string) *Token {
	t := new(Token)
	s := tredis.GetRedis(accessPrefix + access)
	if s == "" {
		return nil
	}
	err := gkit.GetJson(s, t)
	if err != nil {
		return nil
	}
	return t
}

// 设置token数据
func SetTokenAllData(userId int64, key string) error {
	s := tredis.GetRedis(userTokenInfoPrefix + gkit.ToString(userId))
	if s == "" {
		s = key
	} else {
		s += "," + key
	}
	return tredis.SetRedis(userTokenInfoPrefix+gkit.ToString(userId), s, time.Duration(Cache.Set.TokenTime)*time.Second)
}

// 删除所有token数据
func DelTokenAllData(userId int64) {
	s := tredis.GetRedis(userTokenInfoPrefix + gkit.ToString(userId))
	tredis.DeleteRedis(strings.Split(s, ",")...)
	tredis.DeleteRedis(userTokenInfoPrefix + gkit.ToString(userId))
}

// 设置用户角色
func SetRole(userId int64, roles []string) error {
	err := tredis.SetRedis(
		userRolePrefix+gkit.ToString(userId),
		strings.Join(roles, ","),
		time.Duration(Cache.Set.TokenTime)*time.Second)
	// 记录
	_ = SetTokenAllData(userId, userRolePrefix+gkit.ToString(userId))
	return err
}

// 获取用户角色
func GetRole(userId int64) []string {
	s := tredis.GetRedis(userRolePrefix + gkit.ToString(userId))
	return strings.Split(s, ",")
}
