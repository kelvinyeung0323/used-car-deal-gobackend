package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/model"
)

type RoleResRepo struct {
	sqlTemplate *datasource.SqlTemplate
}

func NewRoleResRepo(txMgr *datasource.TransactionManger) *RoleResRepo {
	sqlTemplate, err := datasource.NewSqlTemplate(txMgr, "role-res.sql.tmpl")
	if err != nil {
		panic(fmt.Errorf("初始化角色资源sqlTemplate错误:%v", err))
	}
	return &RoleResRepo{sqlTemplate: sqlTemplate}
}

func (r *RoleResRepo) CreateResOfRole(ctx *gin.Context, roleId string, resIds []string) {

	if resIds == nil || len(resIds) <= 0 {
		log.Warnf("createResOfRole:resIds is empty.roleId:%v", roleId)
		return
	}
	_, err := r.sqlTemplate.Exec(ctx, "createResOfRole", map[string]any{"RoleId": roleId, "ResIds": resIds})
	if err != nil {
		panic(err)
	}
}

func (r *RoleResRepo) DeleteResOfRole(ctx *gin.Context, roleId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteResOfRole", roleId)
	if err != nil {
		panic(err)
	}
}

func (r *RoleResRepo) GetResIdsOfRole(ctx *gin.Context, roleId string) []string {
	var resIds []string
	err := r.sqlTemplate.Select(ctx, &resIds, "getResIdsOfRole", roleId)
	if err != nil {
		panic(err)
	}
	return resIds
}

func (r *RoleResRepo) QueryResOfUser(ctx *gin.Context, form *model.ResQueryForm) []model.Res {
	var resList []model.Res
	err := r.sqlTemplate.Select(ctx, &resList, "queryResOfUser", form)
	if err != nil {
		panic(err)
	}
	return resList

}
