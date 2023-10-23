package repo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/web"
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
	var u model.User
	isExists, err := r.sqlTemplate.Get(ctx, &u, "findById", userId)
	if isExists < 0 {
		panic(err)
	}
	if isExists == 0 {
		return nil
	}
	return &u
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

func (r *UserRepo) Delete(ctx *gin.Context, userId string) {
	_, err := r.sqlTemplate.Exec(ctx, "delete", userId)
	if err != nil {
		panic(err)
	}
}

func (r *UserRepo) PageQuery(ctx *gin.Context, form *model.UserQueryForm) *web.Page[model.User] {
	userList := []model.User{}
	total, err := r.sqlTemplate.PageQuery(ctx, &userList, "query", form)
	if err != nil {
		panic(err)
	}
	return &web.Page[model.User]{Total: total, PageSize: form.PageSize, PageNum: form.PageNum, Data: userList}
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

func (r *UserRepo) EnableUser(ctx *gin.Context, userId string, enabled bool) {
	_, err := r.sqlTemplate.Exec(ctx, "enableUser", map[string]any{"UserId": userId, "Enabled": enabled})
	if err != nil {
		panic(err)
	}
}

func (r *UserRepo) ChangePwd(ctx *gin.Context, form *model.ChangePwdForm) {
	_, err := r.sqlTemplate.Exec(ctx, "changePwd", form)
	if err != nil {
		panic(err)
	}
}
func (r *UserRepo) FilterExistUserIds(ctx *gin.Context, userIds []string) []string {
	ids := []string{}
	err := r.sqlTemplate.Select(ctx, &ids, "filterExistUserIds", userIds)
	if err != nil {
		panic(err)
	}
	return ids
}

func (r *UserRepo) FindAuthCodesOfUser(ctx *gin.Context, userId string) []string {
	authCodes := []string{}
	err := r.sqlTemplate.Select(ctx, &authCodes, "findAuthCodesOfUser", userId)
	if err != nil {
		panic(err)
	}
	return authCodes
}
