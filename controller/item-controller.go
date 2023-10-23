package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
	"used-car-deal-gobackend/base/logger"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

var log = logger.GetInstance()

type ItemController struct {
	itemService *service.ItemService
}

func NewItemController(itemService *service.ItemService) *ItemController {
	return &ItemController{itemService: itemService}
}

// CreateItem 创建商品
func (c *ItemController) CreateItem(ctx *gin.Context) {
	var item model.Item
	if err := ctx.ShouldBindJSON(&item); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.itemService.CreateItem(ctx, &item)
	web.ReturnOK(ctx, nil)
}

// QueryItem 查询Item
func (c *ItemController) QueryItem(ctx *gin.Context) {
	var form model.ItemQueryForm
	if err := ctx.Bind(&form); err != nil {
		log.Warnf("绑定QueryItem参数错误:%v", err)
		web.Err(web.VALID_ERROR)
	}

	page := c.itemService.QueryItems(ctx, &form)
	web.ReturnWithPage(ctx, "", page)
}

func (c *ItemController) UpdateItem(ctx *gin.Context) {

	var form model.Item
	if err := ctx.BindJSON(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	c.itemService.UpdateItem(ctx, &form)
	web.ReturnOK(ctx, nil)
}

func (c *ItemController) DeleteItems(ctx *gin.Context) {
	idStr := ctx.Query("itemIds")
	itemIds := strings.Split(idStr, ",")
	if len(itemIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	c.itemService.DeleteItems(ctx, itemIds)
	web.ReturnOKWithMsg(ctx, nil, "删除商品成功")
}

func (c *ItemController) FindItem(ctx *gin.Context) {
	itemId := ctx.Param("itemId")
	if itemId == "" {
		web.Err(web.VALID_ERROR)
	}
	item := c.itemService.FindItem(ctx, itemId)
	web.ReturnOK(ctx, item)
}
