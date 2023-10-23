package model

import (
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
)

type User struct {
	UserId      string      `json:"userId" db:"user_id"`
	UserName    string      `json:"userName" db:"user_name"`
	Avatar      string      `json:"avatar" db:"avatar"`
	Password    string      `json:"-" db:"password"`
	NickName    string      `json:"nickName" db:"nick_name" `
	Enabled     bool        `json:"enabled" db:"enabled"`
	CreatedAt   *types.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   *types.Time `json:"updatedAt" db:"updated_at"`
	LastLoginAt *types.Time `json:"lastLoginAt" db:"last_login_at"`
	Remark      string      `json:"remark" db:"remark"`
	Roles       []Role      `json:"roles"`
}

// UserQueryForm 用户查询表单
type UserQueryForm struct {
	web.Pageable[User]
	UserName *string     `form:"userName"`
	Enabled  *bool       `form:"enabled"`
	StartAt  *types.Time `form:"startAt"`
	EndAt    *types.Time `form:"endAt"`
}

// UserCreateForm 创建用户表单
type UserForm struct {
	UserId    string   `form:"userId" db:"user_id"`
	UserName  string   `form:"userName" db:"user_name"`
	Avatar    string   `form:"avatar" db:"avatar"`
	Password  string   `form:"password,omitempty" db:"password"`
	NickName  string   `form:"nickName" db:"nick_name"`
	Enabled   bool     `form:"enabled" db:"enabled"`
	Remark    string   `form:"remark" db:"remark"`
	Roles     []Role   `form:"roles"`
	RoleIds   []string `form:"roleIds"`
	AuthCodes []string `form:"authCodes"`
}

type EnableUserForm struct {
	UserId  string `form:"userId" binding:"required"`
	Enabled bool   `form:"enabled" `
}
type ChangePwdForm struct {
	UserId      string `form:"userId"`
	Password    string `form:"password"`
	OldPassword string `form:"oldPassword"`
	NewPassword string `form:"newPassword"`
}
type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
