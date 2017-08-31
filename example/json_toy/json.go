package main

import (
	"errors"
	"time"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/privationel/link"
	"github.com/privationel/link/codec"
)

func init() {
	beego.SetLogFuncCall(true)
}

var channal *link.Channel

func main() {
	channal = link.NewChannel()
	//
	json := codec.Json()
	server, err := link.Listen("tcp", "127.0.0.1:9988", json, 0 /* sync send */)
	if err != nil {
		beego.Error(err)
		return
	}

	go server.Serve()
	server.RegisterRouter("register", RegisterHandle)
	for {
		beego.Info("-----------------------------")
		channal.Fetch(func(session *link.Session) {
			beego.Info(session.ID)
		})
		beego.Info("=============================")
		// beego.Info(channal)
		channal.Fetch(func(session *link.Session) {
			b := HealthCheck(session)
			if b == false {
				channal.DeleteKeyBySession(session)
			}
		})
		time.Sleep(2 * time.Second)

	}
}

type Register struct {
	Ip           string    `json:"ip"`
	RegisterName string    `json:"register_name"`
	NodeType     string    `json:"node_type"`
	RegisterTime time.Time `json:"register_time"`
}

func RegisterHandle(session *link.Session, content interface{}) {
	if _, ok := content.(map[string]interface{}); !ok {
		link.DefalutErrorHandle(session, errors.New("content type error;"))
	}
	b, _ := json.Marshal(content)
	var register_info Register
	err := json.Unmarshal(b, &register_info)
	if err != nil {
		link.DefalutErrorHandle(session, err)
	}
	channal.Put(register_info, session)
	err = session.Send(&link.Response{
		Id:      "register",
		Method:  "register",
		Content: "register success",
	})
	if err != nil {
		beego.Error(err)
	}
	beego.Info(session, "is registered")
}
func HealthCheck(session *link.Session) bool {
	session.Send(&link.Request{
		Method:  "HealthCheck",
		Content: true,
	})
	_, err := session.Receive()
	if err != nil {
		beego.Error("session :", session, "error:", err)
		session.Close()
		return false
	} else {
		return true
	}
}
