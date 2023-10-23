package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

type ParamController struct {
	paramService *service.ParamService
}

func NewParamController(paramService *service.ParamService) *ParamController {
	return &ParamController{paramService: paramService}
}

func (c *ParamController) CreateParam(ctx *gin.Context) {
	var form model.SysParam
	if err := ctx.ShouldBindJSON(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.paramService.CreateParam(ctx, &form)
	web.ReturnOK(ctx, nil)
}

func (c *ParamController) QueryParam(ctx *gin.Context) {
	var form model.SysParamQueryForm
	if err := ctx.Bind(&form); err != nil {
		log.Warnf("绑定QueryItem参数错误:%v", err)
		web.Err(web.VALID_ERROR)
	}

	page := c.paramService.QueryParams(ctx, &form)
	web.ReturnWithPage(ctx, "", page)
}

func (c *ParamController) UpdateParam(ctx *gin.Context) {

	var form model.SysParam
	if err := ctx.BindJSON(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.paramService.UpdateParam(ctx, &form)
	web.ReturnOK(ctx, nil)
}

func (c *ParamController) DeleteParams(ctx *gin.Context) {
	idStr := ctx.Query("paramIds")
	paramIds := strings.Split(idStr, ",")
	if len(paramIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	c.paramService.DeleteParam(ctx, paramIds)
	web.ReturnOKWithMsg(ctx, nil, "删除参数成功")
}

func (c *ParamController) FindParam(ctx *gin.Context) {
	paramId := ctx.Param("paramId")
	if paramId == "" {
		web.Err(web.VALID_ERROR)
	}
	item := c.paramService.FindParam(ctx, paramId)
	web.ReturnOK(ctx, item)
}

//=================条目=========================

func (c *ParamController) CreateItem(ctx *gin.Context) {
	var form model.ParamItem
	if err := ctx.ShouldBindJSON(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.paramService.CreateItem(ctx, &form)
	web.ReturnOK(ctx, nil)
}

func (c *ParamController) UpdateItem(ctx *gin.Context) {

	var form model.ParamItem
	if err := ctx.BindJSON(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.paramService.UpdateItem(ctx, &form)
	web.ReturnOK(ctx, nil)
}

func (c *ParamController) DeleteItems(ctx *gin.Context) {
	idStr := ctx.Query("itemIds")
	itemIds := strings.Split(idStr, ",")
	if len(itemIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	c.paramService.DeleteItem(ctx, itemIds)
	web.ReturnOKWithMsg(ctx, nil, "删除参数条目成功")
}

func (c *ParamController) FindItem(ctx *gin.Context) {
	itemId := ctx.Param("itemId")
	if itemId == "" {
		web.Err(web.VALID_ERROR)
	}
	item := c.paramService.FindItem(ctx, itemId)
	web.ReturnOK(ctx, item)
}
func (c *ParamController) GetItemOfParam(ctx *gin.Context) {
	paramId := ctx.Query("paramId")
	if paramId == "" {
		web.Err(web.VALID_ERROR)
	}

	items := c.paramService.GetItemOfParam(ctx, paramId, "root")
	web.ReturnOK(ctx, items)
}

func (c *ParamController) GetAllParamsWithItems(ctx *gin.Context) {
	items := c.paramService.GetAllParamsWithItems(ctx)
	web.ReturnOK(ctx, items)
}
func (c *ParamController) GetParamWithItemsByKey(ctx *gin.Context) {
	paramKey := ctx.Query("paramKey")
	if paramKey == "" {
		web.Err(web.VALID_ERROR)
	}
	items := c.paramService.GetAllParamsWithItems(ctx)
	web.ReturnOK(ctx, items)
}
