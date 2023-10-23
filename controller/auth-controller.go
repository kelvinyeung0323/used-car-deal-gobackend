package controller

import (
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/security"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

type AuthController struct {
	userService *service.UserService
	authService *service.AuthService
	authHelper  *security.AuthHelper
}

func NewAuthController(userService *service.UserService, authService *service.AuthService, authHelper *security.AuthHelper) *AuthController {
	return &AuthController{userService, authService, authHelper}
}
func (c *AuthController) Login(ctx *gin.Context) {
	var loginForm model.LoginForm
	err := ctx.ShouldBind(&loginForm)
	if err != nil {
		web.Err(err)
	}
	// 校验用户名和密码是否正确
	user := c.authService.FindUserByName(ctx, loginForm.Username)
	if err != nil {
		log.Debugf("查找用户失败%v", err)
		web.Err(web.LOGIN_UNKNOWN)
	}

	if !user.Enabled {
		web.BizErr("用户已停用，请联系管理员")
	}

	if user.Password == *security.EncryptPwd(&loginForm.Password) {
		// 生成Token
		tokenString, _ := security.GenToken(user.UserId, user.UserName)
		//刷新权限到缓存
		ctx.Set("userId", user.UserId)
		auth, err := c.authHelper.RefreshAuthorization(ctx)
		if err != nil {
			log.Debugf("login error:%v", err)
			web.Err(web.FORBIDDEN)
		}
		web.ReturnOK(ctx, gin.H{"token": tokenString, "info": auth})
		log.Debugf("登录成功v")

		return
	}
	log.Debugf("登录失败%v", loginForm)
	web.Err(web.LOGIN_ERROR)
}

func (c *AuthController) CurrentUser(ctx *gin.Context) {
	user, err := c.authHelper.GetAuthorization(ctx)
	if user == nil || err != nil {
		web.BizErr("获取当前用户信息失败.")
	}
	web.ReturnOK(ctx, user)

}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var form model.ChangePwdForm
	err := ctx.ShouldBind(&form)
	if err != nil {
		web.Err(err)
	}
	c.authService.ChangePwd(ctx, &form)
	web.ReturnOKWithMsg(ctx, nil, "修改密码成功")
}

func (c *AuthController) UpdateUserProfile(ctx *gin.Context) {
	var form model.UserForm
	err := ctx.ShouldBind(&form)
	if err != nil {
		web.Err(err)
	}
	c.authService.UpdateUserProfile(ctx, &form)
	web.ReturnOKWithMsg(ctx, nil, "修改个人信息成功")
}
