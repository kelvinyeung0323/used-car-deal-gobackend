package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
)

type ItemCommonRepo struct {
	sqlTemplate *datasource.SqlTemplate
}

func NewItemCommonRepo(txMgr *datasource.TransactionManger) *ItemCommonRepo {
	sqlTemplate, err := datasource.NewSqlTemplate(txMgr, "item-common.sql.tmpl")
	if err != nil {
		panic(fmt.Errorf("初始化商品属性sqlTemplate错误:%v", err))
	}
	return &ItemCommonRepo{sqlTemplate: sqlTemplate}
}

//===============品牌==================================/

func (r *ItemCommonRepo) PageQueryBrands(ctx *gin.Context, form *model.BrandQueryForm) *web.Page[model.Brand] {
	itemList := []model.Brand{}
	total, err := r.sqlTemplate.PageQuery(ctx, &itemList, "queryBrands", form)
	if err != nil {
		panic(err)
	}
	return &web.Page[model.Brand]{Total: total, PageNum: form.PageNum, PageSize: form.PageSize, Data: itemList}
}
func (r *ItemCommonRepo) QueryBrands(ctx *gin.Context, form *model.BrandQueryForm) []model.Brand {
	itemList := []model.Brand{}
	err := r.sqlTemplate.Select(ctx, &itemList, "queryBrands", form)
	if err != nil {
		panic(err)
	}
	return itemList
}
func (r *ItemCommonRepo) FindBrand(ctx *gin.Context, brandId int) *model.Brand {
	u := &model.Brand{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findBrand", brandId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}

func (r *ItemCommonRepo) CreateBrand(ctx *gin.Context, brand *model.Brand) {
	_, err := r.sqlTemplate.Exec(ctx, "createBrand", brand)
	if err != nil {
		panic(err)
	}
}

func (r *ItemCommonRepo) UpdateBrand(ctx *gin.Context, brand *model.Brand) {
	_, err := r.sqlTemplate.Exec(ctx, "updateBrand", brand)
	if err != nil {
		panic(err)
	}
}

func (r *ItemCommonRepo) DeleteBrand(ctx *gin.Context, brandId int) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteBrand", brandId)
	if err != nil {
		panic(err)
	}
}

//========================== 车系  ======================================

func (r *ItemCommonRepo) QuerySeries(ctx *gin.Context, form *model.SeriesQueryForm) *web.Page[model.Series] {
	itemList := []model.Series{}
	total, err := r.sqlTemplate.PageQuery(ctx, &itemList, "querySeries", form)
	if err != nil {
		panic(err)
	}
	return &web.Page[model.Series]{Total: total, PageNum: form.PageNum, PageSize: form.PageSize, Data: itemList}
}
func (r *ItemCommonRepo) FindSeries(ctx *gin.Context, seriesId int) *model.Series {
	u := &model.Series{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findSeries", seriesId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}

func (r *ItemCommonRepo) UpdateSeries(ctx *gin.Context, series *model.Series) {
	_, err := r.sqlTemplate.Exec(ctx, "updateSeries", series)
	if err != nil {
		panic(err)
	}
}

func (r *ItemCommonRepo) DeleteSeries(ctx *gin.Context, seriesId int) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteSeries", seriesId)
	if err != nil {
		panic(err)
	}
}

func (r *ItemCommonRepo) CreateSeries(ctx *gin.Context, series *model.Series) {
	_, err := r.sqlTemplate.Exec(ctx, "createSeries", series)
	if err != nil {
		panic(err)
	}
}

func (r *ItemCommonRepo) GetSeriesOfBrand(ctx *gin.Context, brandId int) []model.Series {
	u := []model.Series{}
	err := r.sqlTemplate.Select(ctx, &u, "getSeriesOfBrand", brandId)
	if err != nil {
		panic(err)
	}
	return u
}
