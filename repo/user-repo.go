package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/model"
)

type UserRepo struct {
	sqlTemplate *datasource.SqlTemplate
}

func NewUserRepo(txMgr *datasource.TransactionManger) *UserRepo {
	sqlTemplate, err := datasource.NewSqlTemplate(txMgr, "user.sql.tmpl")
	if err != nil {
		panic(fmt.Errorf("初始化用户sqlTemplate错误:%v", err))
	}
	return &UserRepo{sqlTemplate: sqlTemplate}
}
func (r *UserRepo) Find(ctx *gin.Context, userId string) *model.User {
	u := &model.User{}
	isExists, err := r.sqlTemplate.Get(ctx, u, "findById", userId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return u
}

func (r *UserRepo) FindByName(ctx *gin.Context, username string) *model.User {
	var u model.User
	isExists, err := r.sqlTemplate.Get(ctx, &u, "findByName", username)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return &u
}

func (r *UserRepo) Update(ctx *gin.Context, user *model.User) {
	_, err := r.sqlTemplate.Exec(ctx, "update", user)
	if err != nil {
		panic(err)
	}
}

func (r *UserRepo) Delete(ctx *gin.Context, userIds []string) {
	_, err := r.sqlTemplate.Exec(ctx, "delete", userIds)
	if err != nil {
		panic(err)
	}
}

func (r *UserRepo) Query(ctx *gin.Context, form *model.UserQueryForm) []model.User {
	var userList []model.User
	err := r.sqlTemplate.Select(ctx, userList, "query", form)
	if err != nil {
		panic(err)
	}
	return userList
}

func (r *UserRepo) Create(ctx *gin.Context, user *model.User) {
	_, err := r.sqlTemplate.Exec(ctx, "create", user)
	if err != nil {
		panic(err)
	}
}

func (r *UserRepo) UpdateLoginTime(ctx *gin.Context) {
	r.sqlTemplate.Exec(ctx, "UpdateLoginTime", nil)
}

func (r *UserRepo) ChangeUserPasswd(ctx *gin.Context) {
	r.sqlTemplate.Exec(ctx, "ChangeUserPasswd", nil)
}
