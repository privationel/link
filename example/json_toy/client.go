package main

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/privationel/link"
	"github.com/privationel/link/codec"
)

func main() {
	json := codec.Json()
	json.Register(link.Request{})
	json.Register(link.Response{})
	client, err := link.Dial("tcp", "127.0.0.1:9988", json, 0)
	if err != nil {
		beego.Error(err)
		return
	}
	clientSessionLoop(client)
	DefaultClientHandle(client)
}

type Register struct {
	Ip           string    `json:"ip"`
	RegisterName string    `json:"register_name"`
	NodeType     string    `json:"node_type"`
	RegisterTime time.Time `json:"register_time"`
}

func clientSessionLoop(session *link.Session) {
	err := session.Send(&link.Request{
		Method: "register",
		Content: &Register{
			Ip:           "xxxxx",
			RegisterName: "test",
			NodeType:     "crawler",
			RegisterTime: time.Now(),
		},
	})
	if err != nil {
		beego.Error(err)
		return
	}

	rsp, err := session.Receive()
	if err != nil {
		beego.Error(err)
		return
	}
	beego.Info(rsp.(*link.Response))
	// err = session.Close()
	// if err != nil {
	// 	beego.Error(err)
	// } else {
	// 	beego.Info("session is closed")
	// }
}

var router *link.Router

func init() {
	router = link.NewRouter()
	router.RouterRegister("HealthCheck", HealthCheck)
}

func DefaultClientHandle(session *link.Session) {
	for {
		beego.Info("in loop")
		rsp, err := session.Receive()
		if err != nil {
			beego.Error(err)
			return
		}
		resp := rsp.(*link.Request)
		f, err := router.GetHandle(resp.Method)
		if err != nil {
			link.DefalutErrorHandle(session, err)
		}
		f(session, resp.Content)
	}
}
func HealthCheck(session *link.Session, content interface{}) {
	beego.Info("CLIENT HEALTH CHECK")
	err := session.Send(&link.Response{
		Method:  "HealthCheck",
		Content: true,
	})
	if err != nil {
		beego.Error(err)
		return
	}
}
