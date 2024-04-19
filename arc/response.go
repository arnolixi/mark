package arc

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Response struct {
	Error error       `json:"error"`
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
}

func (r *Response) response() gin.H {
	var res = gin.H{
		"msg":  r.Msg,
		"code": r.Code,
		"data": r.Data,
	}
	if r.Error != nil {
		res["error"] = r.Error.Error()
	}
	return res
}

func (r *Response) setData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) setMsg(msg string) *Response {
	r.Msg = msg
	return r
}

func (r *Response) JSON(c *gin.Context) {
	c.JSON(http.StatusOK, r.response())
}

type bindFunc func(any) error

func BindExec(f bindFunc, data any) *Response {
	if err := f(data); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return Unwrap(err, nil, "解析失败！")
		}
		return Unwrap(errs, errs.Translate(T()))
	}
	return ResOK()

}

func ResOK() *Response {
	return &Response{
		Msg:  "ok",
		Code: 0,
	}
}

func ResError(err error) *Response {
	return &Response{
		Error: err,
		Code:  CodeERROR,
		Msg:   err.Error(),
	}
}

func Result(args ...interface{}) *Response {
	switch len(args) {
	case 1:
		err, ok := args[0].(error)
		if ok {
			return ResError(err)
		}
		return ResOK()

	case 2:
		err, ok := args[1].(error)
		if ok {
			return Unwrap(err, args[0])
		}
		return Unwrap(nil, args[0])

	case 3:
		err, ok := args[2].(error)
		if ok {
			return Unwrap(err, args[0], args[1])
		}
		return Unwrap(nil, args[0], args[1])

	default:
		return ResOK()
	}
}

func Unwrap(err error, args ...interface{}) *Response {
	res := ResOK()
	if err != nil {
		res = ResError(err)
	}
	switch len(args) {
	case 1:
		res.setData(args[0])

	case 2:
		res.setData(args[0]).setMsg(args[1].(string))
	}
	return res
}
