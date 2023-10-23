package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
)

type UserRoleRepo struct {
	sqlTemplate *datasource.SqlTemplate
}

func NewUserRoleRepo(txMgr *datasource.TransactionManger) *UserRoleRepo {
	sqlTemplate, err := datasource.NewSqlTemplate(txMgr, "user-role.sql.tmpl")
	if err != nil {
		panic(fmt.Errorf("初始化用户角色sqlTemplate错误:%v", err))
	}
	return &UserRoleRepo{sqlTemplate: sqlTemplate}
}

func (r *UserRoleRepo) CreateUserRole(ctx *gin.Context, userId string, roleId string) {
	_, err := r.sqlTemplate.Exec(ctx, "CreateUserRole", map[string]any{"UserId": userId, "RoleId": roleId})
	if err != nil {
		panic(err)
	}
}

func (r *UserRoleRepo) DeleteUserRole(ctx *gin.Context, userId string, roleId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteUserRole", map[string]any{"UserId": userId, "RoleId": roleId})
	if err != nil {
		panic(err)
	}
}

func (r *UserRoleRepo) DeleteUserRoleByRoleId(ctx *gin.Context, roleId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteUserRoleByRoleId", roleId)
	if err != nil {
		panic(err)
	}
}

func (r *UserRoleRepo) QueryUserOfRole(ctx *gin.Context, form *model.RoleUserQueryForm) *web.Page[model.User] {
	userList := []model.User{} //不用 var userList []model.User{},实例化后如果没有数据返回支返回空数组而不是nil
	total, err := r.sqlTemplate.PageQuery(ctx, &userList, "queryUserOfRole", form)
	if err != nil {
		panic(err)
	}
	return &web.Page[model.User]{Total: total, PageNum: form.PageNum, PageSize: form.PageSize, Data: userList}
}

func (r *UserRoleRepo) QueryUserIsNotBelongToRole(ctx *gin.Context, form *model.RoleUserQueryForm) *web.Page[model.User] {
	userList := []model.User{}
	total, err := r.sqlTemplate.PageQuery(ctx, &userList, "queryUserIsNotBelongToRole", form)
	if err != nil {
		panic(err)
	}
	return &web.Page[model.User]{Total: total, PageNum: form.PageNum, PageSize: form.PageSize, Data: userList}
}

func (r *UserRoleRepo) DeleteRolesOfUser(ctx *gin.Context, userId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteRolesOfUser", userId)
	if err != nil {
		panic(err)
	}
}

func (r *UserRoleRepo) CreateRolesOfUser(ctx *gin.Context, userId string, roleIds []string) {
	_, err := r.sqlTemplate.Exec(ctx, "createRolesOfUser", map[string]any{"UserId": userId, "RoleIds": roleIds})
	if err != nil {
		panic(err)
	}
}
func (r *UserRoleRepo) CreateUsersOfRole(ctx *gin.Context, roleId string, userIds []string) {
	_, err := r.sqlTemplate.Exec(ctx, "createUsersOfRole", map[string]any{"RoleId": roleId, "UserIds": userIds})
	if err != nil {
		panic(err)
	}
}

func (r *UserRoleRepo) FindRolesByUserId(ctx *gin.Context, userId string) []model.Role {
	roleList := []model.Role{}
	err := r.sqlTemplate.Select(ctx, &roleList, "findRolesByUserId", userId)
	if err != nil {
		panic(err)
	}
	return roleList
}

func (r *UserRoleRepo) QueryRolesOfUser(ctx *gin.Context, form *model.RoleUserQueryForm) []model.Role {
	roleList := []model.Role{}
	err := r.sqlTemplate.Select(ctx, &roleList, "queryRolesOfUser", form)
	if err != nil {
		panic(err)
	}
	return roleList
}
