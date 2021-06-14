package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponsePaging struct {
	Data   interface{} `json:"data"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Total  int         `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewBaseResponsePaging(data interface{}, limit int, offset int, total int) BaseResponsePaging {
	return BaseResponsePaging{
		Data:   data,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}

func NewBaseResponseScroll(data interface{}, scrollId string) BaseResponseScroll {
	return BaseResponseScroll{
		Items:    data,
		ScrollId: scrollId,
	}
}

func NewResponse(code int, data interface{}) (int, interface{}) {
	return code, gin.H{
		"data": data,
	}
}

func NewOKResponse(data interface{}) (int, interface{}) {
	return http.StatusOK, gin.H{
		"data": data,
	}
}

func NewErrorResponse(code int, msg interface{}) (int, interface{}) {
	return code, gin.H{
		"error": msg,
	}
}

func ServiceUnavailable() (int, interface{}) {
	return http.StatusServiceUnavailable, gin.H{
		"error": http.StatusText(http.StatusServiceUnavailable),
	}
}

func ServiceUnavailableMsg(msg interface{}) (int, interface{}) {
	return http.StatusServiceUnavailable, gin.H{
		"error": msg,
	}
}

func BadRequest() (int, interface{}) {
	return http.StatusBadRequest, gin.H{
		"error": http.StatusText(http.StatusBadRequest),
	}
}

func BadRequestMsg(msg interface{}) (int, interface{}) {
	return http.StatusBadRequest, gin.H{
		"error": msg,
	}
}

func BadRequestErroMsg(msg interface{}) (int, interface{}) {
	return http.StatusBadRequest, msg
}
func OKRequestMsg(msg interface{}) (int, interface{}) {
	return http.StatusOK, msg
}

func NotFound() (int, interface{}) {
	return http.StatusNotFound, gin.H{
		"error": http.StatusText(http.StatusNotFound),
	}
}

func NotFoundMsg(msg interface{}) (int, interface{}) {
	return http.StatusNotFound, gin.H{
		"error": msg,
	}
}

func Forbidden() (int, interface{}) {
	return http.StatusForbidden, gin.H{
		"error": "Do not have permission for the request.",
	}
}

func Unauthorized() (int, interface{}) {
	return http.StatusUnauthorized, gin.H{
		"error": http.StatusText(http.StatusUnauthorized),
	}
}

type BaseResponseScroll struct {
	Items    interface{} `json:"data"`
	ScrollId string      `json:"scroll_id"`
}
