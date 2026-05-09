package hertz

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func NewServer(addr string) *server.Hertz {
	h := server.Default(
		server.WithHostPorts(addr),
	)
	return h
}

func Run(h *server.Hertz) error {
	return h.Run()
}

type Response struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}

func Fail(code int32, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
	}
}
