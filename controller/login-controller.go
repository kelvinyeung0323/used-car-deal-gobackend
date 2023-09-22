package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"used-car-deal-gobackend/base/security"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

type LoginController struct {
	userService *service.UserService
}

func NewLoginController(userService *service.UserService) *LoginController {
	return &LoginController{userService}
}
func (c *LoginController) Login(ctx *gin.Context) {
	var loginForm model.LoginForm
	err := ctx.ShouldBind(&loginForm)
	if err != nil {
		web.Err(err)
	}
	//TODO:从数据库读取用户信息
	// 校验用户名和密码是否正确
	user := c.userService.FindUserByName(ctx, loginForm.Username)
	if err != nil {
		log.Printf("查找用户失败%v", err)
		web.Err(web.LOGIN_UNKNOWN)
	}

	if user.Password == *security.EncryptPwd(&loginForm.Password) {
		// 生成Token
		tokenString, _ := security.GenToken(user.UserId, user.UserName)
		web.ReturnOK(ctx, gin.H{"token": tokenString})
		return
	}
	web.Err(web.LOGIN_ERROR)
}
