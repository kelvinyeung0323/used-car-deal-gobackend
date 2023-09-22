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
	userRepo *repo.UserRepo
	txMgr    *datasource.TransactionManger
}

func NewUserService(userRepo *repo.UserRepo, txMgr *datasource.TransactionManger) *UserService {
	return &UserService{userRepo: userRepo, txMgr: txMgr}
}

//FindUserByName 根据用户名称获取用户
func (s *UserService) FindUserByName(ctx *gin.Context, username string) *model.User {
	user := s.userRepo.FindByName(ctx, username)
	if user == nil {
		web.BizErr("用户不存在!")
	}
	return user
}

// QueryUser 查询用户列表
func (s *UserService) QueryUser(ctx *gin.Context, form *model.UserQueryForm) []model.User {
	userList := s.userRepo.Query(ctx, form)
	return userList
}

func (s *UserService) CreateUser(ctx *gin.Context, form *model.UserCreateForm) {

	//开户事务
	s.txMgr.BeginTx(ctx)

	//判断用户是否存在
	u := s.userRepo.FindByName(ctx, form.UserName)
	if u != nil {
		web.BizErr("用户已存在!")
	}
	//TODO:生成初始密码
	//生成UUID
	uuid := uuid.New()
	user := &model.User{}
	user.UserId = uuid.String()
	user.Password = *security.EncryptPwd(&form.Password)
	user.UserName = form.UserName
	user.NickName = form.NickName
	now := types.Time(time.Now())
	user.CreatedAt = &now
	user.UpdatedAt = &now
	s.userRepo.Create(ctx, user)
	s.txMgr.CommitTx(ctx)
}

func (s *UserService) FindUser(ctx *gin.Context, userId string) *model.User {
	user := s.userRepo.Find(ctx, userId)
	if user == nil {
		web.BizErr("用户不存在")
	}
	return user
}

func (s *UserService) UpdateUser(ctx *gin.Context, user *model.User) {
	s.txMgr.BeginTx(ctx)
	//暂时只能修改nickName
	//1.判断用户是否存在
	u := s.userRepo.Find(ctx, user.UserId)
	if u == nil {
		web.BizErr("用户不存在!")
	}
	s.userRepo.Update(ctx, user)
	s.txMgr.CommitTx(ctx)
}

func (s *UserService) DeleteUser(ctx *gin.Context, userIds []string) {
	s.txMgr.BeginTx(ctx)
	s.userRepo.Delete(ctx, userIds)
	s.txMgr.CommitTx(ctx)
}
