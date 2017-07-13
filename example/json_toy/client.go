package main

import (
	"log"

	"github.com/astaxie/beego"
	"github.com/privationel/link"
	"github.com/privationel/link/codec"
)

func main() {
	json := codec.Json()
	json.Register(link.Request{})
	client, err := link.Dial("tcp", "127.0.0.1:9988", json, 0)
	if err != nil {
		beego.Error(err)
		return
	}
	clientSessionLoop(client)
}
func clientSessionLoop(session *link.Session) {
	for i := 0; i < 10; i++ {
		err := session.Send(&link.Request{
			Method:  "test",
			Content: "content",
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
		log.Printf("Receive: %d", rsp)
	}
}
