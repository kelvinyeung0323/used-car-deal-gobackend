package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/model"
)

type ResRepo struct {
	sqlTemplate *datasource.SqlTemplate
}

func NewResRepo(txMgr *datasource.TransactionManger) *ResRepo {
	sqlTemplate, err := datasource.NewSqlTemplate(txMgr, "res.sql.tmpl")
	if err != nil {
		panic(fmt.Errorf("初始化用户sqlTemplate错误:%v", err))
	}
	return &ResRepo{sqlTemplate: sqlTemplate}
}
func (r *ResRepo) Find(ctx *gin.Context, resId string) *model.Res {
	var u model.Res
	isExists, err := r.sqlTemplate.Get(ctx, &u, "findById", resId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return &u
}

func (r *ResRepo) FindByName(ctx *gin.Context, resName string) *model.Res {
	var u model.Res
	isExists, err := r.sqlTemplate.Get(ctx, &u, "findByName", resName)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return &u
}

func (r *ResRepo) Update(ctx *gin.Context, role *model.Res) {
	_, err := r.sqlTemplate.Exec(ctx, "update", role)
	if err != nil {
		panic(err)
	}
}

func (r *ResRepo) Delete(ctx *gin.Context, resId string) {
	_, err := r.sqlTemplate.Exec(ctx, "delete", resId)
	if err != nil {
		panic(err)
	}
}

func (r *ResRepo) Query(ctx *gin.Context, form *model.ResQueryForm) []model.Res {
	resList := []model.Res{}
	err := r.sqlTemplate.Select(ctx, &resList, "query", form)
	if err != nil {
		panic(err)
	}
	return resList
}

func (r *ResRepo) Create(ctx *gin.Context, res *model.Res) {
	_, err := r.sqlTemplate.Exec(ctx, "create", res)
	if err != nil {
		panic(err)
	}
}

func (r *ResRepo) CheckParentIdIsChild(ctx *gin.Context, resId string, parentId string) bool {
	idList := []string{}
	err := r.sqlTemplate.Select(ctx, &idList, "checkParentIdIsChild", map[string]string{"ParentId": parentId, "ResId": resId})
	if err != nil {
		panic(err)
	}
	if len(idList) > 0 {
		return true
	}
	return false
}

func (r *ResRepo) FilterExistResIds(ctx *gin.Context, resIds []string) []string {
	filteredIds := []string{}
	err := r.sqlTemplate.Select(ctx, &filteredIds, "filterExistResIds", resIds)
	if err != nil {
		panic(err)
	}
	return filteredIds
}

func (r *ResRepo) DeleteRoleResByResId(ctx *gin.Context, resId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteRoleResByResId", resId)
	if err != nil {
		panic(err)
	}
}
