package web

import (
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/logger"
)

var log = logger.GetInstance()

type webError struct {
	Code int64
	Msg  string
}

var (
	// LOGIN_UNKNOWN 错误代码
	LOGIN_UNKNOWN = &webError{202, "用户不存在"}
	LOGIN_ERROR   = &webError{203, "账号或密码错误"}
	TOKEN_ERROR   = &webError{204, "权限验证错误"}
	VALID_ERROR   = &webError{300, "参数错误"}
	ERROR         = &webError{400, "操作失败"}
	UNAUTHORIZED  = &webError{401, "您还未登录"}
	FORBIDDEN     = &webError{403, "访问受限"}
	NOT_FOUND     = &webError{404, "资源不存在"}
	INNER_ERROR   = &webError{500, "系统发生异常"}
	BIZ_FAIL      = &webError{501, "业务错误"}
)

func Err(err any) {
	panic(err)
}

// BizErr 其他业务相关的错误统一使用此方法
func BizErr(msg string) {
	panic(&webError{Code: 501, Msg: msg})
}

// ErrorHandleMiddleware 统一异常处理中间件
func ErrorHandleMiddleware(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case *webError:
				ReturnFail(c, t)
			default:
				ReturnFail(c, INNER_ERROR)
			}
			//重新抛出异常，让logger recover 捕捉
			//c.Abort()
			panic(r)
		}
	}()
	c.Next()
}
