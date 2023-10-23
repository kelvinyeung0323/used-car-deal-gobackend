package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const statusOk = 200

// RestResult Code 为1即为操作成功，其他的为具体错误代码
type RestResult struct {
	Data     any    `json:"data,omitempty"`
	Code     int64  `json:"code,omitempty"`
	Msg      any    `json:"msg,omitempty"`
	PageSize *int64 `json:"pageSize,omitempty"`
	PageNum  *int64 `json:"pageNum,omitempty"`
	Total    *int64 `json:"total,omitempty"`
}

func ReturnResult(c *gin.Context, r *RestResult) {
	c.JSON(http.StatusOK, r)
}

func ReturnFail(c *gin.Context, r *webError) {
	c.JSON(http.StatusOK, &RestResult{
		Code:     r.Code,
		Data:     nil,
		Msg:      r.Msg,
		PageSize: nil,
		PageNum:  nil,
		Total:    nil,
	})
}

func ReturnOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &RestResult{
		Code:     statusOk,
		Data:     data,
		Msg:      nil,
		PageSize: nil,
		PageNum:  nil,
		Total:    nil,
	})

}
func ReturnOKWithMsg(c *gin.Context, data any, msg any) {
	c.JSON(http.StatusOK, &RestResult{
		Code:     statusOk,
		Data:     data,
		Msg:      msg,
		PageSize: nil,
		PageNum:  nil,
		Total:    nil,
	})

}

func ReturnPage(c *gin.Context, msg any, data any, pageSize int64, pageNum int64, total int64) {

	c.JSON(http.StatusOK, &RestResult{
		Code:     statusOk,
		Data:     data,
		Msg:      nil,
		PageSize: &pageSize,
		PageNum:  &pageNum,
		Total:    &total,
	})
}

func ReturnWithPage[T any](c *gin.Context, msg any, page *Page[T]) {

	c.JSON(http.StatusOK, &RestResult{
		Data:     page.Data,
		Code:     statusOk,
		PageSize: &page.PageSize,
		PageNum:  &page.PageNum,
		Total:    &page.Total,
	})
}

func ReturnPageWithMsg(c *gin.Context, data any, pageSize int64, pageNum int64, total int64, msg string) {

	c.JSON(http.StatusOK, &RestResult{
		Code:     statusOk,
		Data:     data,
		Msg:      msg,
		PageSize: &pageSize,
		PageNum:  &pageNum,
		Total:    &total,
	})
}
