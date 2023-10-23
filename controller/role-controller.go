package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

type RoleController struct {
	roleService *service.RoleService
}

func NewRoleController(roleService *service.RoleService) *RoleController {
	return &RoleController{roleService: roleService}
}

func (r *RoleController) Query(ctx *gin.Context) {
	var form model.RoleQueryForm

	if err := ctx.BindQuery(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	rolePage := r.roleService.Query(ctx, &form)
	if rolePage == nil {
		rolePage = &web.Page[model.Role]{}
	}
	web.ReturnWithPage(ctx, "", rolePage)
}

func (r *RoleController) Find(ctx *gin.Context) {
	roleId := ctx.Param("roleId")
	role := r.roleService.Find(ctx, roleId)
	web.ReturnOK(ctx, role)
}

func (r *RoleController) Create(ctx *gin.Context) {
	var form model.Role
	if err := ctx.ShouldBind(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	r.roleService.Create(ctx, &form)
	web.ReturnOKWithMsg(ctx, nil, "创建角色成功")

}

func (r *RoleController) Update(ctx *gin.Context) {
	var form model.Role
	if err := ctx.ShouldBind(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	r.roleService.Update(ctx, &form)
	web.ReturnOKWithMsg(ctx, nil, "修改角色成功")
}

func (r *RoleController) Delete(ctx *gin.Context) {
	roleIdStr := ctx.Query("roleIds")
	roleIds := strings.Split(roleIdStr, ",")
	if len(roleIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	r.roleService.Delete(ctx, roleIds)
	web.ReturnOKWithMsg(ctx, nil, "删除角色成功")
}
func (r *RoleController) EnabledRole(ctx *gin.Context) {
	var form model.EnabledRoleForm
	if err := ctx.ShouldBind(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	r.roleService.EnabledRole(ctx, form.RoleId, form.Enabled)
	web.ReturnOKWithMsg(ctx, nil, "修改角色状态成功")
}

func (r *RoleController) QueryUserOfRole(ctx *gin.Context) {
	var form model.RoleUserQueryForm

	if err := ctx.BindQuery(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	rolePage := r.roleService.QueryUserOfRole(ctx, &form)
	if rolePage == nil {
		rolePage = &web.Page[model.User]{Data: []model.User{}}
	}
	web.ReturnWithPage(ctx, "", rolePage)
}

func (r *RoleController) CreateUsersOfRole(ctx *gin.Context) {

	var usersOfRole model.UsersOfRole
	if err := ctx.ShouldBind(&usersOfRole); err != nil {
		web.Err(web.VALID_ERROR)
	}

	r.roleService.CreateUsersOfRole(ctx, usersOfRole)
	web.ReturnOKWithMsg(ctx, nil, "添加角色用户成功")
}

func (r *RoleController) DeleteUserRole(ctx *gin.Context) {

	var usersOfRole model.UsersOfRole
	if err := ctx.ShouldBind(&usersOfRole); err != nil {
		web.Err(web.VALID_ERROR)
	}

	r.roleService.DeleteUsersOfRole(ctx, usersOfRole)
	web.ReturnOKWithMsg(ctx, nil, "删降角色用户成功")
}
