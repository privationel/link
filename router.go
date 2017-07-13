package link

import (
	"errors"
	"log"

	"github.com/astaxie/beego"
)

type Router struct {
	routerMap map[string]HandleFunction
}
type HandleFunction func(session *Session, request_body interface{})

func NewRouter() (r *Router) {
	r = new(Router)
	r.routerMap = make(map[string]HandleFunction)
	return r
}

func (r *Router) routerRegister(router string, f HandleFunction) {
	r.routerMap[router] = f
}

type Request struct {
	Id      string      `json:"id"`
	Method  string      `json:"method"`
	Content interface{} `json:"content"`
}
type Response struct {
	Id      string      `json:"id"`
	Method  string      `json:"method"`
	Content interface{} `json:"content"`
}

func (s *Server) DefalutHandle(session *Session) {
	req, err := session.Codec().Receive()
	if err != nil {
		log.Println("request content:", req, " error:", err)
		DefalutErrorHandle(session, err)
		return
	}
	value := req.(*Request)
	if value == nil {
		log.Println("request body error or request body is nil.")
		DefalutErrorHandle(session, errors.New("request body error or request body is nil."))
		return
	}
	f := s.router.routerMap[value.Method]
	f(session, value.Content)
}

func DefalutErrorHandle(session *Session, err error) {
	if err != nil {
		errors.New("there is no error message!")
	}
	e := session.Send(&Response{
		Method:  "error",
		Content: err.Error(),
	})
	if e != nil {
		beego.Error(e)
		return
	}
}
