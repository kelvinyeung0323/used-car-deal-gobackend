package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
)

type ParamRepo struct {
	txMgr       *datasource.TransactionManger
	sqlTemplate *datasource.SqlTemplate
}

func NewParamRepo(txMgr *datasource.TransactionManger) *ParamRepo {
	sqlTemplate, err := datasource.NewSqlTemplate(txMgr, "param.sql.tmpl")
	if err != nil {
		panic(fmt.Errorf("初始化系统参数sqlTemplate错误:%v", err))
	}
	return &ParamRepo{txMgr: txMgr, sqlTemplate: sqlTemplate}

}

func (r *ParamRepo) FindParamById(ctx *gin.Context, paramId string) *model.SysParam {
	u := &model.SysParam{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findParam", paramId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}

func (r *ParamRepo) FindParamByKey(ctx *gin.Context, paramKey string) *model.SysParam {
	u := &model.SysParam{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findParamByKey", paramKey)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}

func (r *ParamRepo) UpdateParam(ctx *gin.Context, param *model.SysParam) {
	_, err := r.sqlTemplate.Exec(ctx, "updateParam", param)
	if err != nil {
		panic(err)
	}
}

func (r *ParamRepo) DeleteParam(ctx *gin.Context, paramId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteParam", paramId)
	if err != nil {
		panic(err)
	}
}

func (r *ParamRepo) PageQueryParams(ctx *gin.Context, form *model.SysParamQueryForm) *web.Page[model.SysParam] {
	params := []model.SysParam{}
	total, err := r.sqlTemplate.PageQuery(ctx, &params, "queryParams", form)
	if err != nil {
		panic(err)
	}
	return &web.Page[model.SysParam]{Total: total, PageNum: form.PageNum, PageSize: form.PageSize, Data: params}
}
func (r *ParamRepo) QueryParams(ctx *gin.Context, form *model.SysParamQueryForm) []model.SysParam {
	params := []model.SysParam{}
	err := r.sqlTemplate.Select(ctx, &params, "queryParams", form)
	if err != nil {
		panic(err)
	}
	return params
}

func (r *ParamRepo) CreateParam(ctx *gin.Context, param *model.SysParam) {
	_, err := r.sqlTemplate.Exec(ctx, "createParam", param)
	if err != nil {
		panic(err)
	}
}

// ====================条目=======================================
func (r *ParamRepo) CreateItem(ctx *gin.Context, item *model.ParamItem) {
	_, err := r.sqlTemplate.Exec(ctx, "createItem", item)
	if err != nil {
		panic(err)
	}
}

func (r *ParamRepo) UpdateItem(ctx *gin.Context, item *model.ParamItem) {
	_, err := r.sqlTemplate.Exec(ctx, "updateItem", item)
	if err != nil {
		panic(err)
	}
}

func (r *ParamRepo) FindItem(ctx *gin.Context, itemId string) *model.ParamItem {
	u := &model.ParamItem{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findItem", itemId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}

func (r *ParamRepo) FindItemsByParamId(ctx *gin.Context, paramId string) []model.ParamItem {
	u := []model.ParamItem{}
	err := r.sqlTemplate.Select(ctx, &u, "findItemsByParamId", paramId)
	if err != nil {
		panic(err)
	}
	return u
}

func (r *ParamRepo) FindItemsByParentId(ctx *gin.Context, parentId string) []model.ParamItem {
	u := []model.ParamItem{}
	err := r.sqlTemplate.Select(ctx, &u, "findItemsByParentId", parentId)
	if err != nil {
		panic(err)
	}
	return u
}

func (r *ParamRepo) FindItemsOfParamByParentId(ctx *gin.Context, paramId string, parentId string) []model.ParamItem {
	u := []model.ParamItem{}
	err := r.sqlTemplate.Select(ctx, &u, "findItemsOfParamByParentId", map[string]string{"ParamId": paramId, "ParentId": parentId})
	if err != nil {
		panic(err)
	}
	return u
}

func (r *ParamRepo) DeleteItem(ctx *gin.Context, itemId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteItem", itemId)
	if err != nil {
		panic(err)
	}
}

func (r *ParamRepo) DeleteItemByParamId(ctx *gin.Context, paramId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteItemByParamId", paramId)
	if err != nil {
		panic(err)
	}
}
