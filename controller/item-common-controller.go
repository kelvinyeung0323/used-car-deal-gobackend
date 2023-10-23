package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

type ItemCommonController struct {
	commonService *service.ItemCommonService
}

func NewItemCommonController(commonService *service.ItemCommonService) *ItemCommonController {
	return &ItemCommonController{commonService: commonService}
}

func (c *ItemCommonController) CreateBrand(ctx *gin.Context) {
	var brand model.Brand
	if err := ctx.ShouldBindJSON(&brand); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.commonService.CreateBrand(ctx, &brand)
	web.ReturnOK(ctx, nil)
}

func (c *ItemCommonController) QueryBrands(ctx *gin.Context) {
	var form model.BrandQueryForm
	if err := ctx.Bind(&form); err != nil {
		log.Warnf("绑定QueryBrands参数错误:%v", err)
		web.Err(web.VALID_ERROR)
	}

	page := c.commonService.QueryBrand(ctx, &form)
	web.ReturnWithPage(ctx, "", page)
}

func (c *ItemCommonController) UpdateBrand(ctx *gin.Context) {

	var form model.Brand
	if err := ctx.BindJSON(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.commonService.UpdateBrand(ctx, &form)
	web.ReturnOK(ctx, nil)
}

func (c *ItemCommonController) DeleteBrands(ctx *gin.Context) {
	idStr := ctx.Query("brandIds")
	brandIds := strings.Split(idStr, ",")
	var intIds []int
	for _, s := range brandIds {
		i, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			log.Warnf("解释参数错误！%v", err)
			web.Err(web.VALID_ERROR)
		}
		intIds = append(intIds, int(i))
	}
	if len(brandIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	c.commonService.DeleteBrands(ctx, intIds)
	web.ReturnOKWithMsg(ctx, nil, "删除商品成功")
}

func (c *ItemCommonController) FindBrand(ctx *gin.Context) {
	brandId := ctx.Param("brandId")
	if brandId == "" {
		web.Err(web.VALID_ERROR)
	}

	intId, err := strconv.ParseInt(brandId, 0, 64)
	if err != nil {
		log.Warnf("解释参数错误！%v", err)
		web.Err(web.VALID_ERROR)
	}
	item := c.commonService.FindBrand(ctx, int(intId))
	web.ReturnOK(ctx, item)
}

func (c *ItemCommonController) CreateSeries(ctx *gin.Context) {
	var series model.Series
	if err := ctx.ShouldBindJSON(&series); err != nil {
		log.Warnf("绑定CreateSeries参数错误:%v", err)
		web.Err(web.VALID_ERROR)
	}
	c.commonService.CreateSeries(ctx, &series)
	web.ReturnOK(ctx, nil)
}

func (c *ItemCommonController) QuerySeries(ctx *gin.Context) {
	var form model.SeriesQueryForm
	if err := ctx.Bind(&form); err != nil {
		log.Warnf("绑定QueryBrands参数错误:%v", err)
		web.Err(web.VALID_ERROR)
	}

	page := c.commonService.QuerySeries(ctx, &form)
	web.ReturnWithPage(ctx, "", page)
}

func (c *ItemCommonController) UpdateSeries(ctx *gin.Context) {

	var form model.Series
	if err := ctx.BindJSON(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.commonService.UpdateSeries(ctx, &form)
	web.ReturnOK(ctx, nil)
}

func (c *ItemCommonController) DeleteSeries(ctx *gin.Context) {
	str := ctx.Query("seriesIds")
	idsArr := strings.Split(str, ",")
	var intIds []int
	for _, s := range idsArr {
		i, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			log.Warnf("解释参数错误！%v", err)
			web.Err(web.VALID_ERROR)
		}
		intIds = append(intIds, int(i))
	}
	if len(intIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	c.commonService.DeleteSeries(ctx, intIds)
	web.ReturnOKWithMsg(ctx, nil, "删除车系成功")
}

func (c *ItemCommonController) FindSeries(ctx *gin.Context) {
	seriesId := ctx.Param("seriesId")
	if seriesId == "" {
		web.Err(web.VALID_ERROR)
	}

	intId, err := strconv.ParseInt(seriesId, 0, 64)
	if err != nil {
		log.Warnf("解释参数错误！%v", err)
		web.Err(web.VALID_ERROR)
	}
	item := c.commonService.FindSeries(ctx, int(intId))
	web.ReturnOK(ctx, item)
}

func (c *ItemCommonController) GetAllBrandsWithSeries(ctx *gin.Context) {
	brands := c.commonService.GetAllBrandsWithSeries(ctx)
	web.ReturnOK(ctx, brands)
}
