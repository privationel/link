package main

import (
	"github.com/astaxie/beego"
	"github.com/privationel/link"
	"github.com/privationel/link/codec"
)

func main() {
	json := codec.Json()
	server, err := link.Listen("tcp", "127.0.0.1:9988", json, 0 /* sync send */)
	if err != nil {
		beego.Error(err)
		return
	}
	server.RegisterRouter("test", TestHandle)
	server.Serve()
}

func TestHandle(session *link.Session, content interface{}) {
	beego.Info(content)
}
