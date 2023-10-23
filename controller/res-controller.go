package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

type ResController struct {
	resService *service.ResService
}

func NewResController(resService *service.ResService) *ResController {
	return &ResController{resService: resService}
}
func (r *ResController) Query(ctx *gin.Context) {
	var form model.ResQueryForm

	if err := ctx.BindQuery(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	resList := r.resService.Query(ctx, &form)
	if resList == nil {
		resList = []model.Res{}
	}
	web.ReturnOK(ctx, resList)
}

func (r *ResController) Find(ctx *gin.Context) {
	resId := ctx.Param("resId")
	res := r.resService.Find(ctx, resId)
	web.ReturnOK(ctx, res)

}

func (r *ResController) Create(ctx *gin.Context) {
	var form model.Res
	if err := ctx.ShouldBind(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	r.resService.Create(ctx, &form)
	web.ReturnOKWithMsg(ctx, nil, "创建菜单成功")
}

func (r *ResController) Update(ctx *gin.Context) {
	var form model.Res
	if err := ctx.ShouldBind(&form); err != nil {
		web.Err(web.VALID_ERROR)
	}
	r.resService.Update(ctx, &form)
	web.ReturnOKWithMsg(ctx, nil, "修改资源成功")
}

func (r *ResController) Delete(ctx *gin.Context) {
	resIdsStr := ctx.Query("resId")
	resIds := strings.Split(resIdsStr, ",")
	if len(resIds) <= 0 {
		web.Err(web.VALID_ERROR)
	}
	r.resService.Delete(ctx, resIds)
	web.ReturnOKWithMsg(ctx, nil, "删除资源成功")
}
