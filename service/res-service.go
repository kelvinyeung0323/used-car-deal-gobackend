package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
	"time"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/repo"
)

type ResService struct {
	resRepo *repo.ResRepo
	txMgr   *datasource.TransactionManger
}

func NewResService(resRepo *repo.ResRepo, txMgr *datasource.TransactionManger) *ResService {
	return &ResService{resRepo: resRepo, txMgr: txMgr}
}

// Query 资源查询
// 返回结果为资源树 不分页
func (s *ResService) Query(ctx *gin.Context, form *model.ResQueryForm) []model.Res {
	if form.ResType != nil || form.ResName != "" {
		//返回列表
		return s.resRepo.Query(ctx, form)

	} else {
		//返回资源树
		if form.ParentId == "" {
			form.ParentId = "root"
		}
		root := s.resRepo.Query(ctx, form)
		for i, p := range root {
			var f1 = &model.ResQueryForm{ParentId: p.ResId, Sortable: form.Sortable}
			c := s.Query(ctx, f1)
			//range创建副本，如果c.Children=c会不生效
			root[i].Children = c
		}

		return root
	}

}

// Create 创建
func (s *ResService) Create(ctx *gin.Context, res *model.Res) {
	if res.ResName == "" {
		web.BizErr("菜单名称不能为空")
	}
	if res.ParentId == "" {
		res.ParentId = "root"
	}
	if res.ParentId != "root" {
		parent := s.resRepo.Find(ctx, res.ParentId)
		if parent == nil {
			web.BizErr("父菜单不存在！")
		}
	}
	if res.AuthCode != "" {
		match := authCodeRegex.MatchString(res.AuthCode)
		if !match {
			web.BizErr("权限标识只能是大小写字母、数字,冒号\":\"、下划线\"_\"或点\".\"")
		}
	}
	res.ResId = strings.ReplaceAll(uuid.New().String(), "-", "")
	var now types.Time = types.Time(time.Now())
	res.CreatedAt = &now
	res.UpdatedAt = &now
	s.resRepo.Create(ctx, res)
}

func (s *ResService) Find(ctx *gin.Context, resId string) *model.Res {
	return s.resRepo.Find(ctx, resId)
}

// Update 更新
func (s *ResService) Update(ctx *gin.Context, res *model.Res) {
	if res.ResName == "" {
		web.BizErr("资源名称不能为空")
	}

	//当前菜单下的所有节点都不能作为此节点的父节点
	if s.resRepo.CheckParentIdIsChild(ctx, res.ResId, res.ParentId) {
		web.BizErr("当前菜单及其下的所有节点都不能作为此节点的父节点")
	}
	if res.AuthCode != "" {
		match := authCodeRegex.MatchString(res.AuthCode)
		if !match {
			web.BizErr("权限标识只能是大小写字母、数字,冒号\":\"、下划线\"_\"或点\".\"")
		}
	}
	now := types.Time(time.Now())
	res.UpdatedAt = &now
	s.resRepo.Update(ctx, res)
}

// Delete 删除资源，支持删除多个资源
func (s *ResService) Delete(ctx *gin.Context, resIds []string) {
	s.txMgr.BeginTx(ctx)
	//把子资源也删除
	for _, resId := range resIds {
		var f1 = &model.ResQueryForm{ParentId: resId}
		children := s.Query(ctx, f1)
		var childIds []string
		for _, child := range children {
			childIds = append(childIds, child.ResId)
			s.Delete(ctx, childIds)
		}
		s.resRepo.Delete(ctx, resId)
		//把角色资源关联关系也删掉
		s.resRepo.DeleteRoleResByResId(ctx, resId)
	}

	s.txMgr.CommitTx(ctx)
}
