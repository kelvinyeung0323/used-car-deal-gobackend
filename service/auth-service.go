package service

import (
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/logger"
	"used-car-deal-gobackend/base/security"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/repo"
)

var log = logger.GetInstance()

type AuthService struct {
	txMgr        *datasource.TransactionManger
	userRepo     *repo.UserRepo
	resRepo      *repo.ResRepo
	roleRepo     *repo.RoleRepo
	userRoleRepo *repo.UserRoleRepo
	roleResRepo  *repo.RoleResRepo
}

func NewAuthService(txMgr *datasource.TransactionManger, userRepo *repo.UserRepo, roleRepo *repo.RoleRepo, resRepo *repo.ResRepo, userRoleRepo *repo.UserRoleRepo, roleResRepo *repo.RoleResRepo) *AuthService {
	return &AuthService{txMgr: txMgr, userRepo: userRepo, roleRepo: roleRepo, resRepo: resRepo, userRoleRepo: userRoleRepo, roleResRepo: roleResRepo}
}

func (s *AuthService) GetAuthorization(ctx *gin.Context) *security.Authorization {
	i, ok := ctx.Get("userId")
	if !ok {
		web.Err(web.UNAUTHORIZED)
	}
	userId := i.(string)
	user := s.userRepo.Find(ctx, userId)
	if user == nil {
		log.Warnf("授权服务找不到对应的用户信息!")
		web.Err(web.UNAUTHORIZED)
	}
	auth := &security.Authorization{User: user}
	enabled := true
	resQueryForm := &model.ResQueryForm{Enabled: &enabled}
	resQueryForm.Sorts = "sort.asc"
	roleQueryForm := &model.RoleUserQueryForm{}

	isAdmin := false
	if user.UserName == "admin" {
		isAdmin = true
	} else {
		resQueryForm.UserId = userId
		roleQueryForm.UserId = userId
	}
	//获取菜单authCodes
	resList := s.roleResRepo.QueryResOfUser(ctx, resQueryForm)
	//获取角色authCodes
	roles := s.userRoleRepo.QueryRolesOfUser(ctx, roleQueryForm)
	auth.Roles = roles
	//获取菜单树
	auth.Res = s.GetAuthResTree(ctx, isAdmin, resQueryForm)
	for _, r := range roles {
		if r.Enabled {
			auth.RoleAuthCodes = append(auth.RoleAuthCodes, r.AuthCode)
		}
	}
	auth.Permissions = map[string]bool{}
	for _, res := range resList {
		if res.AuthCode != "" {
			auth.Permissions[res.AuthCode] = true
		}
	}

	return auth
}

func (s *AuthService) GetAuthKey(ctx *gin.Context) string {
	userId, ok := ctx.Get("userId")
	if !ok {
		web.Err(web.UNAUTHORIZED)
	}
	return userId.(string)
}

// FindUserByName 根据用户名称获取用户
func (s *AuthService) FindUserByName(ctx *gin.Context, username string) *model.User {
	user := s.userRepo.FindByName(ctx, username)
	if user == nil {
		web.BizErr("用户不存在!")
	}
	return user
}

func (s *AuthService) GetAuthResTree(ctx *gin.Context, isAdmin bool, form *model.ResQueryForm) []model.Res {

	//返回资源树,查询分配的资源
	if form.ParentId == "" {
		form.ParentId = "root"
	}
	var root []model.Res
	if isAdmin {
		root = s.resRepo.Query(ctx, form)
	} else {
		root = s.roleResRepo.QueryResOfUser(ctx, form)
	}

	for i, p := range root {
		var f1 = &model.ResQueryForm{UserId: form.UserId, Enabled: form.Enabled, ParentId: p.ResId, Sortable: form.Sortable}
		c := s.GetAuthResTree(ctx, isAdmin, f1)
		//range创建副本，如果c.Children=c会不生效
		root[i].Children = c
	}

	return root
}

func (s *AuthService) ChangePwd(ctx *gin.Context, form *model.ChangePwdForm) {
	s.txMgr.BeginTx(ctx)
	userId := s.GetAuthKey(ctx)
	if userId == "" {
		web.BizErr("获取不到当前用户信息")
	}
	user := s.userRepo.Find(ctx, userId)
	if user == nil {
		web.BizErr("获取不到当前用户信息")
	}
	encryptPwd := security.EncryptPwd(&form.OldPassword)
	if *encryptPwd != user.Password {
		web.BizErr("旧密码不正确")
	}
	form.UserId = userId
	form.Password = *security.EncryptPwd(&form.NewPassword)
	s.userRepo.ChangePwd(ctx, form)

	s.txMgr.CommitTx(ctx)
}

func (s *AuthService) UpdateUserProfile(ctx *gin.Context, form *model.UserForm) {
	s.txMgr.BeginTx(ctx)
	userId := s.GetAuthKey(ctx)
	if userId == "" {
		web.BizErr("获取不到当前用户信息")
	}
	user := s.userRepo.Find(ctx, userId)
	if user == nil {
		web.BizErr("获取不到当前用户信息")
	}
	user.NickName = form.NickName
	s.userRepo.Update(ctx, user)
	s.txMgr.CommitTx(ctx)
}
