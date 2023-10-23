package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

type UserController struct {
	userService *service.UserService
	fileService *service.FileService
}

func NewUserController(userService *service.UserService, fileService *service.FileService) *UserController {
	return &UserController{userService: userService}
}

func (s *UserController) Query(ctx *gin.Context) {
	form := model.UserQueryForm{}
	if err := ctx.BindQuery(&form); err != nil {
		log.Debugf("error:%v", err)
		web.Err(web.VALID_ERROR)
	}
	page := s.userService.QueryUser(ctx, &form)
	web.ReturnWithPage[model.User](ctx, nil, page)
}
func (c *UserController) Find(ctx *gin.Context) {
	userId := ctx.Param("userId")
	user := c.userService.FindUser(ctx, userId)
	//转换图片路径
	user.Avatar = c.fileService.TransToUrl(user.Avatar)
	web.ReturnOK(ctx, user)

}
func (c *UserController) Update(ctx *gin.Context) {
	var form model.UserForm
	if err := ctx.ShouldBind(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.userService.UpdateUser(ctx, &form)
	web.ReturnOKWithMsg(ctx, nil, "修改用户成功")
}

func (c *UserController) Delete(ctx *gin.Context) {
	userIdStr := ctx.Query("userIds")
	userIds := strings.Split(userIdStr, ",")
	if len(userIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	c.userService.DeleteUser(ctx, userIds)
	web.ReturnOKWithMsg(ctx, nil, "删除用户成功")
}

func (s *UserController) Create(ctx *gin.Context) {
	var form model.UserForm
	if err := ctx.ShouldBind(&form); err != nil {
		log.Debugf("error:%v", err)
		web.Err(web.VALID_ERROR)
	}
	s.userService.CreateUser(ctx, &form)

	web.ReturnOKWithMsg(ctx, nil, "用户创建成功！")
}

func (s *UserController) EnableUser(ctx *gin.Context) {
	var form model.EnableUserForm
	if err := ctx.ShouldBind(&form); err != nil {
		log.Debugf("error:%v", err)
		web.Err(web.VALID_ERROR)
	}
	s.userService.EnableUser(ctx, form.UserId, form.Enabled)
	web.ReturnOKWithMsg(ctx, nil, "用户状态修改成功！")
}

func (s *UserController) ChangePwd(ctx *gin.Context) {
	var form model.ChangePwdForm
	if err := ctx.ShouldBind(&form); err != nil {
		log.Debugf("error:%v", err)
		web.Err(web.VALID_ERROR)
	}
	s.userService.ChangePwd(ctx, &form)
	web.ReturnOKWithMsg(ctx, nil, "用户密码修改成功！")
}
