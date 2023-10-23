package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
)

type ItemRepo struct {
	sqlTemplate *datasource.SqlTemplate
}

func NewItemRepo(txMgr *datasource.TransactionManger) *ItemRepo {
	sqlTemplate, err := datasource.NewSqlTemplate(txMgr, "item.sql.tmpl")
	if err != nil {
		panic(fmt.Errorf("初始化商品sqlTemplate错误:%v", err))
	}
	return &ItemRepo{sqlTemplate: sqlTemplate}
}

func (r *ItemRepo) FindItemById(ctx *gin.Context, itemId string) *model.Item {
	u := &model.Item{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findItemById", itemId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}

func (r *ItemRepo) UpdateItem(ctx *gin.Context, item *model.Item) {
	_, err := r.sqlTemplate.Exec(ctx, "updateItem", item)
	if err != nil {
		panic(err)
	}
}

func (r *ItemRepo) DeleteItem(ctx *gin.Context, itemId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteItem", itemId)
	if err != nil {
		panic(err)
	}
}

func (r *ItemRepo) QueryItems(ctx *gin.Context, form *model.ItemQueryForm) *web.Page[model.Item] {
	itemList := []model.Item{}
	total, err := r.sqlTemplate.PageQuery(ctx, &itemList, "queryItems", form)
	if err != nil {
		panic(err)
	}
	return &web.Page[model.Item]{Total: total, PageNum: form.PageNum, PageSize: form.PageSize, Data: itemList}
}

func (r *ItemRepo) CreateItem(ctx *gin.Context, item *model.Item) {
	_, err := r.sqlTemplate.Exec(ctx, "createItem", item)
	if err != nil {
		panic(err)
	}
}

func (r *ItemRepo) CreateItemMedia(ctx *gin.Context, media *model.ItemMedia) {
	_, err := r.sqlTemplate.Exec(ctx, "createItemMedia", media)
	if err != nil {
		panic(err)
	}
}

func (r *ItemRepo) GetMediaOfItem(ctx *gin.Context, itemId string) []model.ItemMedia {
	medias := []model.ItemMedia{}
	err := r.sqlTemplate.Select(ctx, &medias, "getMediaOfItem", itemId)
	if err != nil {
		panic(err)
	}
	return medias
}

func (r *ItemRepo) DeleteMediaOfItem(ctx *gin.Context, itemId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteMediaOfItem", itemId)
	if err != nil {
		panic(err)
	}
}

func (r *ItemRepo) DeleteItemMedia(ctx *gin.Context, media *model.ItemMedia) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteItemMedia", media)
	if err != nil {
		panic(err)
	}
}

func (r *ItemRepo) FindItemMediaById(ctx *gin.Context, mediaId string) *model.ItemMedia {
	u := &model.ItemMedia{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findItemMediaById", mediaId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}
