package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"regexp"
	"strings"
	"time"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/repo"
)

var authCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9._:%]+$`)

type RoleService struct {
	roleRepo     *repo.RoleRepo
	resRepo      *repo.ResRepo
	userRepo     *repo.UserRepo
	userRoleRepo *repo.UserRoleRepo
	roleResRepo  *repo.RoleResRepo
	txMgr        *datasource.TransactionManger
}

func NewRoleService(userRepo *repo.UserRepo, roleRepo *repo.RoleRepo, resRespo *repo.ResRepo, userRoleRepo *repo.UserRoleRepo, roleResRepo *repo.RoleResRepo, txMgr *datasource.TransactionManger) *RoleService {
	return &RoleService{userRepo: userRepo, roleRepo: roleRepo, resRepo: resRespo, userRoleRepo: userRoleRepo, roleResRepo: roleResRepo, txMgr: txMgr}
}

func (s *RoleService) Query(ctx *gin.Context, form *model.RoleQueryForm) *web.Page[model.Role] {
	return s.roleRepo.PageQuery(ctx, form)
}

func (s *RoleService) Find(ctx *gin.Context, roleId string) *model.Role {
	role := s.roleRepo.Find(ctx, roleId)
	if role != nil {
		resIds := s.roleResRepo.GetResIdsOfRole(ctx, roleId)
		role.ResIds = resIds
	}
	return role
}
func (s *RoleService) Update(ctx *gin.Context, form *model.Role) *model.Role {
	s.txMgr.BeginTx(ctx)
	if strings.Trim(form.RoleId, " ") == "" {
		web.BizErr("角色ID不能为空")
	}
	if strings.Trim(form.RoleName, " ") == "" {
		web.BizErr("角色名称不能为空")
	}
	if strings.Trim(form.AuthCode, " ") == "" {
		web.BizErr("权限标识不能为空")
	}
	match := authCodeRegex.MatchString(form.AuthCode)
	if !match {
		web.BizErr("权限标识只能是大小写字母、数字,冒号\":\"、下划线\"_\"或点\".\"")
	}
	isAuthCodeExists := s.roleRepo.CheckAuthCodeExists(ctx, form)
	if isAuthCodeExists {
		web.BizErr("权限标识已存在，请使用其他标识")
	}

	now := types.Time(time.Now())
	form.UpdatedAt = &now
	s.roleRepo.Update(ctx, form)
	s.roleResRepo.DeleteResOfRole(ctx, form.RoleId)
	if len(form.ResIds) > 0 {
		resIds := s.resRepo.FilterExistResIds(ctx, form.ResIds)
		if len(resIds) > 0 {
			s.roleResRepo.CreateResOfRole(ctx, form.RoleId, resIds)
		}
	}
	s.txMgr.CommitTx(ctx)
	return nil
}

func (s *RoleService) Create(ctx *gin.Context, form *model.Role) *model.Role {
	s.txMgr.BeginTx(ctx)
	if strings.Trim(form.RoleName, " ") == "" {
		web.BizErr("角色名称不能为空")
	}
	if strings.Trim(form.AuthCode, " ") == "" {
		web.BizErr("权限标识不能为空")
	}

	match := authCodeRegex.MatchString(form.AuthCode)
	if !match {
		web.BizErr("权限标识只能是大小写字母、数字,冒号\":\"、下划线\"_\"或点\".\"")
	}

	isAuthCodeExists := s.roleRepo.CheckAuthCodeExists(ctx, form)
	if isAuthCodeExists {
		web.BizErr("权限标识已存在，请使用其他标识")
	}
	roleId := strings.ReplaceAll(uuid.New().String(), "-", "")
	form.RoleId = roleId
	now := types.Time(time.Now())
	form.CreatedAt = &now
	form.UpdatedAt = &now
	s.roleRepo.Create(ctx, form)
	s.roleResRepo.DeleteResOfRole(ctx, form.RoleId)
	if len(form.ResIds) > 0 {
		resIds := s.resRepo.FilterExistResIds(ctx, form.ResIds)
		if len(resIds) > 0 {
			s.roleResRepo.CreateResOfRole(ctx, form.RoleId, resIds)
		}
	}
	s.txMgr.CommitTx(ctx)
	return nil
}

func (s *RoleService) Delete(ctx *gin.Context, roleIds []string) {
	s.txMgr.BeginTx(ctx)
	for _, roleId := range roleIds {
		s.roleRepo.Delete(ctx, roleId)
		//把用户角色产系也一起删掉
		s.userRoleRepo.DeleteUserRoleByRoleId(ctx, roleId)
	}
	s.txMgr.CommitTx(ctx)
}

func (s *RoleService) EnabledRole(ctx *gin.Context, roleId string, enabled bool) {
	s.txMgr.BeginTx(ctx)
	if roleId == "" {
		web.BizErr("角色ID不能为空")
	}
	s.roleRepo.EnableRole(ctx, roleId, enabled)
	s.txMgr.CommitTx(ctx)
}

func (s *RoleService) QueryUserOfRole(ctx *gin.Context, form *model.RoleUserQueryForm) *web.Page[model.User] {
	//IsBelongToRole标识，为true表示查询属于角色的用户，false 表示查询不属于角色的用户
	if form.IsBelongToRole {
		return s.userRoleRepo.QueryUserOfRole(ctx, form)
	} else {
		return s.userRoleRepo.QueryUserIsNotBelongToRole(ctx, form)
	}
}

func (s *RoleService) CreateUsersOfRole(ctx *gin.Context, ur model.UsersOfRole) {
	s.txMgr.BeginTx(ctx)
	if ur.RoleId == "" {
		web.BizErr("角色ID不能为空")
	}
	if len(ur.UserIds) <= 0 {
		web.BizErr("用户ID不能为空")
	}
	//判断ID是否存在
	role := s.roleRepo.Find(ctx, ur.RoleId)
	if role == nil {
		web.Err("角色不存在！")
	}
	filteredIds := s.userRepo.FilterExistUserIds(ctx, ur.UserIds)
	if len(filteredIds) <= 0 {
		web.Err("用户不存在")
	}

	s.userRoleRepo.CreateUsersOfRole(ctx, ur.RoleId, filteredIds)
	s.txMgr.CommitTx(ctx)
}

func (s *RoleService) DeleteUsersOfRole(ctx *gin.Context, ur model.UsersOfRole) {
	s.txMgr.BeginTx(ctx)
	if ur.RoleId == "" || ur.UserIds == nil || len(ur.UserIds) <= 0 {
		web.BizErr("角色ID或用户id不能为空")
	}
	for _, userId := range ur.UserIds {
		s.userRoleRepo.DeleteUserRole(ctx, userId, ur.RoleId)
	}
	s.txMgr.CommitTx(ctx)
}
