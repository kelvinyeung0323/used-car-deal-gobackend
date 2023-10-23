package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

import (
	"strings"
	"used-car-deal-gobackend/base/web"
)

type ItemSpecController struct {
	specService *service.ItemSpecService
	fileService *service.FileService
}

func NewItemSpecController(specService *service.ItemSpecService, fileService *service.FileService) *ItemSpecController {
	return &ItemSpecController{specService: specService}
}

// CreateItem 创建商品
func (c *ItemSpecController) CreateSpec(ctx *gin.Context) {
	var item model.ItemSpec
	if err := ctx.ShouldBindJSON(&item); err != nil {
		log.Errorf("绑定CreateSpec参数错误:%v", err)
		web.Err(web.VALID_ERROR)
	}
	c.specService.CreateSpec(ctx, &item)
	web.ReturnOK(ctx, nil)
}

// QuerySpecs 查询规格信息
func (c *ItemSpecController) QuerySpecs(ctx *gin.Context) {
	var form model.ItemSpecQueryForm
	if err := ctx.Bind(&form); err != nil {
		log.Errorf("绑定QuerySpecs参数错误:%v", err)
		web.Err(web.VALID_ERROR)
	}
	if form.BrandIds != nil && len(form.BrandIds) > 0 {
		form.BrandIds = strings.Split(form.BrandIds[0], ",")
	}
	if form.SeriesIds != nil && len(form.SeriesIds) > 0 {
		form.SeriesIds = strings.Split(form.SeriesIds[0], ",")
	}
	page := c.specService.QuerySpecs(ctx, &form)
	web.ReturnWithPage(ctx, "", page)
}

func (c *ItemSpecController) UpdateSpec(ctx *gin.Context) {

	var form model.ItemSpec
	if err := ctx.BindJSON(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.specService.UpdateSpec(ctx, &form)
	web.ReturnOK(ctx, nil)
}

func (c *ItemSpecController) DeleteSpecs(ctx *gin.Context) {
	idStr := ctx.Query("itemIds")
	itemIds := strings.Split(idStr, ",")
	if len(itemIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	c.specService.DeleteSpecs(ctx, itemIds)
	web.ReturnOKWithMsg(ctx, nil, "删除商品规格成功")
}

func (c *ItemSpecController) FindSpec(ctx *gin.Context) {
	specId := ctx.Param("specId")
	if specId == "" {
		web.Err(web.VALID_ERROR)
	}
	spec := c.specService.FindSpec(ctx, specId)
	web.ReturnOK(ctx, spec)
}

func (c *ItemSpecController) GetColorsOfSpec(ctx *gin.Context) {
	specId := ctx.Query("specId")
	if specId == "" {
		web.Err(web.VALID_ERROR)
	}
	specs := c.specService.GetColorsOfSpec(ctx, specId)
	web.ReturnOK(ctx, specs)
}

func (c *ItemSpecController) CreateSpecColor(ctx *gin.Context) {
	var item model.ItemSpecColor
	if err := ctx.ShouldBindJSON(&item); err != nil {
		log.Errorf("绑定CCreateSpecColor参数错误:%v", err)
		web.Err(web.VALID_ERROR)
	}
	c.specService.CreateSpecColor(ctx, &item)
	web.ReturnOK(ctx, nil)
}

func (c *ItemSpecController) DeleteSpecColor(ctx *gin.Context) {
	idStr := ctx.Query("colorIds")
	colorIds := strings.Split(idStr, ",")
	if len(colorIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	c.specService.DeleteSpecColor(ctx, colorIds)
	web.ReturnOKWithMsg(ctx, nil, "删除商品规格颜色成功")
}

func (c *ItemSpecController) FindSpecColor(ctx *gin.Context) {
	colorId := ctx.Param("colorId")
	if colorId == "" {
		web.Err(web.VALID_ERROR)
	}
	color := c.specService.FindSpecColor(ctx, colorId)
	web.ReturnOK(ctx, color)
}

func (c *ItemSpecController) UpdateSpecColor(ctx *gin.Context) {
	var form model.ItemSpecColor
	if err := ctx.BindJSON(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.specService.UpdateSpecColor(ctx, &form)
	web.ReturnOK(ctx, nil)
}

func (c *ItemSpecController) UploadSpecMedia(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	log.Debug("upload file:", file.Filename)
	if err != nil {
		log.Warnf("文件上传失败,%v", err)
		web.Err(web.ERROR)
	}
	specId := ctx.Query("specId")
	colorId := ctx.Query("colorId")
	catalogStr := ctx.Query("catalog")
	if specId == "" || colorId == "" || catalogStr == "" {
		log.Warnf("图片上传参数错误，specId:%v,colorId:%v,catalog:%v", specId, colorId, catalogStr)
		web.Err(web.VALID_ERROR)
	}
	intCatalog, err := strconv.ParseInt(catalogStr, 0, 64)
	if err != nil {
		log.Warnf("图片类目转换错误：%v", catalogStr)
	}
	catalog := model.ItemSpecMediaCatalog(intCatalog)

	//区分文件类型图片放到Images/文件放到Files
	fileName, filePath, err := c.fileService.SaveToTmp(ctx, file)
	if err != nil {
		log.Warnf("文件保存到临时文件夹失败:%v", err)
		web.Err(web.ERROR)
	}
	filePath, err = c.fileService.CopyToImages(ctx, filePath)
	if err != nil {
		log.Warnf("文件复制到目标文件夹失败:%v", err)
		web.Err(web.ERROR)
	}
	//TODO:判断媒体类型
	media := &model.ItemSpecMedia{SpecId: specId, ColorId: colorId, Catalog: catalog, Location: filePath, MediaType: model.MediaImage}
	c.specService.UploadSpecMedia(ctx, media)
	uf := &model.UploadedFile{Id: media.MediaId, Name: fileName, Url: c.fileService.TransToUrl(filePath), Size: file.Size}
	web.ReturnOK(ctx, uf)
}

func (c *ItemSpecController) DeleteSpecMedia(ctx *gin.Context) {
	idStr := ctx.Query("mediaIds")
	MediaIds := strings.Split(idStr, ",")
	if len(MediaIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	c.specService.DeleteSpecMedia(ctx, MediaIds)
	web.ReturnOKWithMsg(ctx, nil, "删除商品规格图片成功")
}
