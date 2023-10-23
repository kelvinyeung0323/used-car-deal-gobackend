package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/logger"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
)

type RoleRepo struct {
	sqlTemplate *datasource.SqlTemplate
}

var log = logger.GetInstance()

func NewRoleRepo(txMgr *datasource.TransactionManger) *RoleRepo {
	sqlTemplate, err := datasource.NewSqlTemplate(txMgr, "role.sql.tmpl")
	if err != nil {
		panic(fmt.Errorf("初始化用户sqlTemplate错误:%v", err))
	}
	return &RoleRepo{sqlTemplate: sqlTemplate}
}
func (r *RoleRepo) Find(ctx *gin.Context, roleId string) *model.Role {
	u := &model.Role{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findById", roleId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}

func (r *RoleRepo) Update(ctx *gin.Context, role *model.Role) {
	_, err := r.sqlTemplate.Exec(ctx, "update", role)
	if err != nil {
		panic(err)
	}
}

func (r *RoleRepo) Delete(ctx *gin.Context, roleId string) {
	_, err := r.sqlTemplate.Exec(ctx, "delete", roleId)
	if err != nil {
		panic(err)
	}
}

func (r *RoleRepo) PageQuery(ctx *gin.Context, form *model.RoleQueryForm) *web.Page[model.Role] {
	var roleList []model.Role
	total, err := r.sqlTemplate.PageQuery(ctx, &roleList, "query", form)
	if err != nil {
		panic(err)
	}
	return &web.Page[model.Role]{Total: total, PageNum: form.PageNum, PageSize: form.PageSize, Data: roleList}
}

func (r *RoleRepo) Query(ctx *gin.Context, form *model.RoleQueryForm) []model.Role {
	var roleList []model.Role
	err := r.sqlTemplate.Select(ctx, &roleList, "query", form)
	if err != nil {
		panic(err)
	}
	return roleList
}

func (r *RoleRepo) Create(ctx *gin.Context, role *model.Role) {
	_, err := r.sqlTemplate.Exec(ctx, "create", role)
	if err != nil {
		panic(err)
	}
}

func (r *RoleRepo) CheckAuthCodeExists(ctx *gin.Context, role *model.Role) bool {
	var roleList []model.Role
	err := r.sqlTemplate.Select(ctx, &roleList, "checkAuthCodeExists", role)
	if err != nil {
		panic(err)
	}
	if len(roleList) > 0 {
		return true
	}
	return false
}

func (r *RoleRepo) FilterExistRoleIds(ctx *gin.Context, roleIds []string) []string {
	var filteredIds []string
	err := r.sqlTemplate.Select(ctx, &filteredIds, "filterExistRoleIds", roleIds)
	if err != nil {
		panic(err)
	}
	return filteredIds
}

func (r *RoleRepo) EnableRole(ctx *gin.Context, roleId string, enabled bool) {
	_, err := r.sqlTemplate.Exec(ctx, "enableRole", map[string]any{"RoleId": roleId, "Enabled": enabled})
	if err != nil {
		panic(err)
	}
}
