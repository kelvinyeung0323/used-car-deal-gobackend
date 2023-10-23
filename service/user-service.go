package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/security"
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/repo"
)

type UserService struct {
	userRepo     *repo.UserRepo
	roleRepo     *repo.RoleRepo
	userRoleRepo *repo.UserRoleRepo
	txMgr        *datasource.TransactionManger
	fileService  *FileService
}

func NewUserService(userRepo *repo.UserRepo, roleRepo *repo.RoleRepo, userRoleRepo *repo.UserRoleRepo, fileService *FileService, txMgr *datasource.TransactionManger) *UserService {
	return &UserService{userRepo: userRepo, roleRepo: roleRepo, userRoleRepo: userRoleRepo, fileService: fileService, txMgr: txMgr}
}

// QueryUser 查询用户列表
func (s *UserService) QueryUser(ctx *gin.Context, form *model.UserQueryForm) *web.Page[model.User] {
	userPage := s.userRepo.PageQuery(ctx, form)
	if userPage != nil && userPage.Data != nil && len(userPage.Data) > 0 {
		for i := 0; i < len(userPage.Data); i++ {
			user := &userPage.Data[i]
			roles := s.userRoleRepo.FindRolesByUserId(ctx, user.UserId)
			user.Roles = roles
		}
	}
	return userPage
}

func (s *UserService) CreateUser(ctx *gin.Context, form *model.UserForm) {

	//开户事务
	s.txMgr.BeginTx(ctx)

	if form.UserName == "" {
		web.BizErr("用户名称不能为空")
	}
	//判断用户是否存在
	u := s.userRepo.FindByName(ctx, form.UserName)
	if u != nil {
		web.BizErr("用户已存在!")
	}
	if form.Password == "" {
		web.BizErr("密码不能为空")
	}

	//生成UUID
	newUserId := uuid.New()
	user := &model.User{}
	user.UserId = newUserId.String()
	user.NickName = form.NickName
	user.UserName = form.UserName
	user.Password = *security.EncryptPwd(&form.Password)
	user.Remark = form.Remark
	user.Enabled = form.Enabled
	now := types.Time(time.Now())
	user.CreatedAt = &now
	user.UpdatedAt = &now

	//头像处理，将图片归档到头像目录
	if form.Avatar != "" {
		filePath, err := s.fileService.MoveToAvatarByUrl(ctx, form.Avatar)
		log.Warnf("保存头像失败！%v", err)
		if err != nil {
			web.Err(web.INNER_ERROR)
		}
		user.Avatar = filePath
	}

	s.userRepo.Create(ctx, user)

	//添加角色
	if form.RoleIds != nil && len(form.RoleIds) > 0 {
		//检查角色是否存在
		roleIds := s.roleRepo.FilterExistRoleIds(ctx, form.RoleIds)
		if len(roleIds) > 0 {
			s.userRoleRepo.CreateRolesOfUser(ctx, user.UserId, form.RoleIds)
		}
	}

	s.txMgr.CommitTx(ctx)
}

func (s *UserService) EnableUser(ctx *gin.Context, userId string, enabled bool) {
	s.txMgr.BeginTx(ctx)
	if userId == "" {
		web.BizErr("用户ID不能为空")
	}
	s.userRepo.EnableUser(ctx, userId, enabled)
	s.txMgr.CommitTx(ctx)
}

func (s *UserService) ChangePwd(ctx *gin.Context, form *model.ChangePwdForm) {
	s.txMgr.BeginTx(ctx)
	if form.UserId == "" && form.Password == "" {
		web.BizErr("用户ID和密码不能为空")
	}
	form.Password = *security.EncryptPwd(&form.Password)
	s.userRepo.ChangePwd(ctx, form)
	s.txMgr.CommitTx(ctx)
}

func (s *UserService) FindUser(ctx *gin.Context, userId string) *model.User {
	user := s.userRepo.Find(ctx, userId)
	if user == nil {
		web.BizErr("用户不存在")
	}
	//获取角色ID
	roles := s.userRoleRepo.FindRolesByUserId(ctx, user.UserId)
	user.Roles = roles
	return user
}

func (s *UserService) UpdateUser(ctx *gin.Context, form *model.UserForm) {
	s.txMgr.BeginTx(ctx)
	//暂时只能修改nickName
	//1.判断用户是否存在
	u := s.userRepo.Find(ctx, form.UserId)
	if u == nil {
		web.BizErr("用户不存在!")
	}
	user := &model.User{}
	user.UserId = form.UserId
	user.UserName = form.UserName
	user.NickName = form.NickName
	user.Enabled = form.Enabled
	user.Remark = form.Remark

	if user.Password != "" {
		user.Password = *security.EncryptPwd(&form.Password)
	}

	//头像处理，将图片归档到头像目录
	if form.Avatar != "" {
		if u.Avatar != "" && u.Avatar != s.fileService.TransToFilePath(form.Avatar) {
			//先删除旧的图片
			err := s.fileService.RemoveFileByPath(u.Avatar)
			if err != nil {
				log.Warnf("删除旧头像失败%v", err)
			}
		}
		filePath, err := s.fileService.MoveToAvatarByUrl(ctx, form.Avatar)
		if err != nil {
			log.Warnf("保存头像失败！%v", err)
			web.Err(web.INNER_ERROR)
		}
		user.Avatar = filePath
	}

	now := types.Time(time.Now())
	user.UpdatedAt = &now
	s.userRepo.Update(ctx, user)

	//先删除后添加
	s.userRoleRepo.DeleteRolesOfUser(ctx, user.UserId)
	//添加角色
	if form.RoleIds != nil && len(form.RoleIds) > 0 {
		//检查角色是否存在
		roleIds := s.roleRepo.FilterExistRoleIds(ctx, form.RoleIds)
		if len(roleIds) > 0 {
			s.userRoleRepo.CreateRolesOfUser(ctx, user.UserId, form.RoleIds)
		}
	}

	s.txMgr.CommitTx(ctx)
}

func (s *UserService) DeleteUser(ctx *gin.Context, userIds []string) {
	s.txMgr.BeginTx(ctx)
	for _, userId := range userIds {
		s.userRepo.Delete(ctx, userId)
	}
	s.txMgr.CommitTx(ctx)
	//TODO:删除头像图片
}
