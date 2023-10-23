package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/repo"
)

type ParamService struct {
	paramCache *cache.Cache
	paramRepo  *repo.ParamRepo
	txMgr      *datasource.TransactionManger
}

func NewParamService(paramRepo *repo.ParamRepo, txMgr *datasource.TransactionManger) *ParamService {
	paramCache := cache.New(24*time.Hour, 24*time.Hour)
	return &ParamService{paramRepo: paramRepo, paramCache: paramCache, txMgr: txMgr}
}

func (s *ParamService) CreateParam(ctx *gin.Context, form *model.SysParam) {
	s.txMgr.BeginTx(ctx)
	//2006-01-02 15:04:05
	if form.ParamKey == "" || form.ParamName == "" {
		web.BizErr("名称和主键称不参为空")
	}

	form.ParamId = strings.ReplaceAll(uuid.New().String(), "-", "")
	now := types.Time(time.Now())
	form.CreatedAt = &now
	form.UpdatedAt = &now

	s.paramRepo.CreateParam(ctx, form)
	s.txMgr.CommitTx(ctx)
}

func (s *ParamService) UpdateParam(ctx *gin.Context, form *model.SysParam) {
	s.txMgr.BeginTx(ctx)
	//2006-01-02 15:04:05
	if form.ParamId == "" {
		web.BizErr("id 不能为空")
	}
	if form.ParamKey == "" || form.ParamName == "" {
		web.BizErr("名称和主键称不参为空")
	}

	now := types.Time(time.Now())
	form.UpdatedAt = &now
	s.paramRepo.UpdateParam(ctx, form)
	s.txMgr.CommitTx(ctx)
}

func (s *ParamService) DeleteParam(ctx *gin.Context, paramIds []string) {
	s.txMgr.BeginTx(ctx)
	for _, paramId := range paramIds {
		s.paramRepo.DeleteParam(ctx, paramId)

	}
	s.txMgr.CommitTx(ctx)
}
func (s *ParamService) QueryParams(ctx *gin.Context, form *model.SysParamQueryForm) *web.Page[model.SysParam] {
	page := s.paramRepo.PageQueryParams(ctx, form)
	return page
}
func (s *ParamService) FindParam(ctx *gin.Context, paramId string) *model.SysParam {
	param := s.paramRepo.FindParamById(ctx, paramId)
	return param
}

//===============条目======================================

func (s *ParamService) CreateItem(ctx *gin.Context, form *model.ParamItem) {
	s.txMgr.BeginTx(ctx)
	//2006-01-02 15:04:05
	if form.ParamId == "" {
		web.BizErr("参数Id不能为空")
	}
	param := s.paramRepo.FindParamById(ctx, form.ParamId)
	if param == nil {
		web.BizErr("参数不存在")
	}
	if form.ParentId == "" || form.ParentId == "root" {
		form.ParentId = "root"
	} else {
		if p := s.paramRepo.FindItem(ctx, form.ParentId); p == nil {
			web.BizErr("父条目不存在")
		}
	}
	form.ItemId = strings.ReplaceAll(uuid.New().String(), "-", "")
	now := types.Time(time.Now())
	form.CreatedAt = &now
	form.UpdatedAt = &now
	s.paramRepo.CreateItem(ctx, form)
	s.txMgr.CommitTx(ctx)
}

func (s *ParamService) UpdateItem(ctx *gin.Context, form *model.ParamItem) {
	s.txMgr.BeginTx(ctx)
	//2006-01-02 15:04:05
	if form.ParamId == "" {
		web.BizErr("参数ID不能为空")
	}
	if form.ItemId == "" {
		web.BizErr("条目ID不能为空")
	}
	param := s.paramRepo.FindParamById(ctx, form.ParamId)
	if param == nil {
		web.BizErr("参数不存在")
	}
	if form.ParentId == "" || form.ParentId == "root" {
		form.ParentId = "root"
	} else {
		if p := s.paramRepo.FindItem(ctx, form.ParentId); p == nil {
			web.BizErr("父条目不存在")
		}
	}
	//不做item是否存在数据库的判断
	now := types.Time(time.Now())
	form.UpdatedAt = &now
	s.paramRepo.UpdateItem(ctx, form)
	s.txMgr.CommitTx(ctx)
}

func (s *ParamService) DeleteItem(ctx *gin.Context, itemIds []string) {
	s.txMgr.BeginTx(ctx)
	for _, itemId := range itemIds {
		//删除其子条目
		children := s.paramRepo.FindItemsByParentId(ctx, itemId)
		childrenIds := []string{}
		for _, child := range children {
			childrenIds = append(childrenIds, child.ItemId)
		}
		s.DeleteItem(ctx, childrenIds)
		s.paramRepo.DeleteItem(ctx, itemId)

	}
	s.txMgr.CommitTx(ctx)
}

func (s *ParamService) FindItem(ctx *gin.Context, itemId string) *model.ParamItem {
	item := s.paramRepo.FindItem(ctx, itemId)
	return item
}

func (s *ParamService) GetItemOfParam(ctx *gin.Context, paramId string, parentId string) []model.ParamItem {

	items := s.paramRepo.FindItemsOfParamByParentId(ctx, paramId, parentId)
	for i := 0; i < len(items); i++ {
		item := &items[i]
		item.Children = s.GetItemOfParam(ctx, paramId, item.ItemId)

	}
	return items
}

func (s *ParamService) GetItemOfParamByKey(ctx *gin.Context, key string) {

}

func (s *ParamService) GetAllParamsWithItems(ctx *gin.Context) []model.SysParam {
	params := s.paramRepo.QueryParams(ctx, &model.SysParamQueryForm{})
	for i := 0; i < len(params); i++ {
		param := &params[i]
		items := s.GetItemOfParam(ctx, param.ParamId, "root")
		param.Items = items
	}
	return params
}

func (s *ParamService) GetParamWithItemsByKey(ctx *gin.Context, paramKey string) *model.SysParam {
	param := s.paramRepo.FindParamByKey(ctx, paramKey)
	if param == nil {
		web.BizErr("找不到对应的参数")
	}
	items := s.GetItemOfParam(ctx, param.ParamId, "root")
	param.Items = items
	return param
}
