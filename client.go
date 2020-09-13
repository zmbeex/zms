package zms

import (
	"github.com/gorilla/websocket"
	"github.com/zmbeex/gkit"
	"net"
	"os"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

type Client struct {
	Conn        *websocket.Conn
	Url         string   `title:"ws的url"`
	Reconnect   chan int `title:"重新连接信道"`
	isReconnect bool     //正在重新连接
}

func (c *Client) reconnect() {
	if c.isReconnect {
		return
	}
	gkit.Debug("开始连接注册中心")
	c.isReconnect = true
	defer func() {
		c.isReconnect = false
	}()

	u := new(Zurl)
	timeStamp := time.Now().Unix()
	signString := "zms.test|" + gkit.ToString(timeStamp) + "|" + Cache.Uuid
	signString += "|" + Cache.Set.UserName + "|" + Cache.Set.Password
	u.Add(Cache.Set.GatewayHost + "/gateway?")
	u.Push("timeStamp", timeStamp)
	u.Push("uid", Cache.Uuid)
	u.Push("sign", gkit.GetSHA(signString))
	u.Push("userName", Cache.Set.UserName)
	c.Url = u.String()
	gkit.Debug(c.Url)
	conn, _, err := websocket.DefaultDialer.Dial(c.Url, nil)
	if err != nil {
		gkit.Error("连接失败 ------>")
		gkit.Error(signString)
		gkit.Error(gkit.GetSHA(signString))
	}
	if err != nil {
		gkit.Error("连接失败:" + err.Error())
		time.Sleep(5 * time.Second)
		c.reconnect()
		return
	}
	gkit.Info("连接注册中心成功")
	c.Conn = conn
}

// 处理心跳
func (c *Client) Hearbeat() {
	for {
		if c.Conn == nil {
			time.Sleep(1 * time.Second)
			c.reconnect()
			continue
		}
		time.Sleep(time.Second * 10)
		result := new(Result)
		result.Code = "zms.system.heart.beats"
		result.Uuid = Cache.Uuid
		result.Status = 1
		c.Write(result)
	}
}

//处理websocket消息
func (c *Client) handleMessage() {
	for {
		if c.Conn == nil {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	for {
		_, message, err := (*c.Conn).ReadMessage()
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				gkit.Warn("网络错误,连接中断")
				time.Sleep(10 * time.Second)
				c.reconnect()
				continue
			}
			if _, ok := err.(*websocket.CloseError); ok {
				gkit.Info("服务端关闭了连接")
				time.Sleep(10 * time.Second)
				c.reconnect()
				continue
			}

			gkit.Warn("错误类型" + reflect.TypeOf(err).String())
			gkit.Debug("错误：" + err.Error())
			os.Exit(1)
			return
		}
		go c.dealMessage(string(message))
	}
}

// 处理结果
func (c *Client) dealMessage(msg string) {
	gkit.Debug("接收到信息：" + msg)
	params := new(Params)
	result := new(Result)
	defer func() {
		r := recover()
		if r != nil {
			gkit.Info("报文异常===========================")
			gkit.Error(string(debug.Stack()))
		}
		c.Write(result)
	}()

	if msg == "zms.system.api.doc" {
		result.Code = msg
		for _, val := range Cache.ServerInfo {
			result.List = append(result.List, val)
		}
		c.Write(result)
		return
	}
	err := gkit.GetJson(msg, params)
	if err != nil {
		gkit.Error("报文数据异常：" + params.Code)
		result.Status = -1
		return
	}
	params.Code = strings.TrimPrefix(params.Code, Cache.Set.UserName+"#")
	result.Code = params.Code
	result.Uuid = params.Uuid
	r := Cache.ServerRouterMap[params.Code]
	if r == nil {
		result.Status = -2
		gkit.Error("收到无效服务：" + params.Code)
		return
	}
	z := InitZms(params, result, r)
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err, ok := r.(error)
		if ok {
			z.SetStatus(-2)
			z.SetNote(err.Error())
			return
		}
		n, ok := r.(int)
		if ok {
			z.SetStatus(n * -1)
			return
		}
		note, ok := r.(string)
		if ok {
			z.SetStatus(-3)
			z.SetNote(note)
			return
		}
		z.SetStatus(-4)
		z.SetNote("未知异常")
	}()
	r.Func(z)
	if result.Status == 0 {
		result.Status = 1
	}
}

// 写入数据
func (c *Client) Write(msg *Result) {
	if msg.Status == 0 {
		msg.Status = 1
	}
	gkit.Debug("发送报文zms：" + gkit.SetJson(msg))
	err := c.Conn.WriteJSON(msg)
	gkit.CheckError(err, "数据发送失败")
}
