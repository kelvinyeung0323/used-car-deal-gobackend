package repo

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
)

type ItemSpecRepo struct {
	sqlTemplate *datasource.SqlTemplate
}

func NewItemSpecRepo(txMgr *datasource.TransactionManger) *ItemSpecRepo {
	sqlTemplate, err := datasource.NewSqlTemplate(txMgr, "item-spec.sql.tmpl")
	if err != nil {
		panic(fmt.Errorf("初始化商品sqlTemplate错误:%v", err))
	}
	return &ItemSpecRepo{sqlTemplate: sqlTemplate}
}

func (r *ItemSpecRepo) FindSpecById(ctx *gin.Context, specId string) *model.ItemSpec {
	u := &model.ItemSpec{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findSpecById", specId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}

	err = json.Unmarshal([]byte(u.ParamsStr), &u.Params)
	if err != nil {
		log.Error("spec params 字段string:[%v] 转struct错误:%v ", u.ParamsStr, err)
	}
	return u
}

func (r *ItemSpecRepo) UpdateSpec(ctx *gin.Context, spec *model.ItemSpec) {
	_, err := r.sqlTemplate.Exec(ctx, "updateSpec", spec)
	if err != nil {
		panic(err)
	}
}

func (r *ItemSpecRepo) DeleteSpec(ctx *gin.Context, specId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteSpec", specId)
	if err != nil {
		panic(err)
	}
}

func (r *ItemSpecRepo) QuerySpecs(ctx *gin.Context, form *model.ItemSpecQueryForm) *web.Page[model.ItemSpec] {
	itemList := []model.ItemSpec{}
	total, err := r.sqlTemplate.PageQuery(ctx, &itemList, "querySpecs", form)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(itemList); i++ {
		item := &itemList[i]
		//TODO:后续改为通用方法
		err = json.Unmarshal([]byte(item.ParamsStr), &item.Params)
		if err != nil {
			log.Error("spec params 字段string:[%v] 转struct错误:%v ", item.ParamsStr, err)
		}
	}
	return &web.Page[model.ItemSpec]{Total: total, PageNum: form.PageNum, PageSize: form.PageSize, Data: itemList}
}

func (r *ItemSpecRepo) CreateSpec(ctx *gin.Context, spec *model.ItemSpec) {
	_, err := r.sqlTemplate.Exec(ctx, "createSpec", spec)
	if err != nil {
		panic(err)
	}
}

func (r *ItemSpecRepo) CreateSpecMedia(ctx *gin.Context, media *model.ItemSpecMedia) {
	_, err := r.sqlTemplate.Exec(ctx, "createSpecMedia", media)
	if err != nil {
		panic(err)
	}
}

func (r *ItemSpecRepo) GetMediaOfSpec(ctx *gin.Context, specId string) []model.ItemSpecMedia {
	medias := []model.ItemSpecMedia{}
	err := r.sqlTemplate.Select(ctx, &medias, "getMediaOfSpec", specId)
	if err != nil {
		panic(err)
	}
	return medias
}

func (r *ItemSpecRepo) GetMediaOfColor(ctx *gin.Context, colorId string) []model.ItemSpecMedia {
	medias := []model.ItemSpecMedia{}
	err := r.sqlTemplate.Select(ctx, &medias, "getMediaOfColor", colorId)
	if err != nil {
		panic(err)
	}
	return medias
}

func (r *ItemSpecRepo) DeleteMediaOfSpec(ctx *gin.Context, specId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteMediaOfSpec", specId)
	if err != nil {
		panic(err)
	}
}

func (r *ItemSpecRepo) DeleteSpecMedia(ctx *gin.Context, mediaId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteSpecMedia", mediaId)
	if err != nil {
		panic(err)
	}
}

func (r *ItemSpecRepo) FindSpecMediaById(ctx *gin.Context, mediaId string) *model.ItemSpecMedia {
	u := &model.ItemSpecMedia{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findSpecMediaById", mediaId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}

func (r *ItemSpecRepo) CreateSpecColor(ctx *gin.Context, color *model.ItemSpecColor) {
	_, err := r.sqlTemplate.Exec(ctx, "createSpecColor", color)
	if err != nil {
		panic(err)
	}
}
func (r *ItemSpecRepo) UpdateSpecColor(ctx *gin.Context, color *model.ItemSpecColor) {
	_, err := r.sqlTemplate.Exec(ctx, "updateSpecColor", color)
	if err != nil {
		panic(err)
	}
}

func (r *ItemSpecRepo) GetColorOfSpec(ctx *gin.Context, specId string) []model.ItemSpecColor {
	medias := []model.ItemSpecColor{}
	err := r.sqlTemplate.Select(ctx, &medias, "getColorOfSpec", specId)
	if err != nil {
		panic(err)
	}
	return medias
}

func (r *ItemSpecRepo) DeleteColorOfSpec(ctx *gin.Context, specId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteColorOfSpec", specId)
	if err != nil {
		panic(err)
	}
}

func (r *ItemSpecRepo) DeleteSpecColor(ctx *gin.Context, colorId string) {
	_, err := r.sqlTemplate.Exec(ctx, "deleteSpecColor", colorId)
	if err != nil {
		panic(err)
	}
}

func (r *ItemSpecRepo) FindSpecColorById(ctx *gin.Context, colorId string) *model.ItemSpecColor {
	u := &model.ItemSpecColor{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findSpecColorById", colorId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}
