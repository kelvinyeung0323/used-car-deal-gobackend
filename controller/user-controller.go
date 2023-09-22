package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (*UserController) Query(ctx *gin.Context) {

}
func (*UserController) Find(ctx *gin.Context) {

}
func (*UserController) Update(ctx *gin.Context) {

}
func (*UserController) Delete(ctx *gin.Context) {

}

func (s *UserController) Create(ctx *gin.Context) {
	var form model.UserCreateForm
	if err := ctx.ShouldBind(&form); err != nil {
		log.Printf("error:%v", err)
		web.Err(web.VALID_ERROR)
	}
	s.userService.CreateUser(ctx, &form)
	web.ReturnOKWithMsg(ctx, nil, "用户创建成功！")
}
