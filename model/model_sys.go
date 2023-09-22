package model

import "used-car-deal-gobackend/base/types"

type User struct {
	UserId      string      `json:"userId" db:"user_id"`
	UserName    string      `json:"userName" db:"user_name"`
	Password    string      `json:"password,omitempty" db:"password"`
	NickName    string      `json:"nickName" db:"nick_name" `
	CreatedAt   *types.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   *types.Time `json:"updatedAt" db:"updated_at"`
	LastLoginAt *types.Time `json:"lastLoginAt" db:"last_login_at"`
	RoleList    []Role
	ResList     []Res
}

type Role struct {
	RoleId    string      `json:"roleId" db:"role_id"`
	RoleName  string      `json:"roleName" db:"role_name"`
	RoleCode  string      `json:"roleCode" db:"role_code"`
	CreatedAt *types.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *types.Time `json:"updatedAt" db:"updated_at"`
	ResList   []Res
}

type Res struct {
	ResId     string      `json:"resId" db:"res_id"`
	ResCode   string      `json:"resCode" db:"res_code"`
	ResName   string      `json:"resName" db:"res_name"`
	ResType   int         `json:"resType" db:"res_type"`
	ResUrl    string      `json:"resUrl" db:"res_url"`
	ParentId  string      `json:"parentId" db:"parent_id"`
	CreatedAt *types.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *types.Time `json:"updatedAt" db:"updated_at"`
}

//UserQueryForm 用户查询表单
type UserQueryForm struct {
	Username *string `form:"username"`
}

//UserCreateForm 创建用户表单
type UserCreateForm struct {
	UserName string `json:"username" db:"user_name" binding:"required"`
	Password string `json:"password,omitempty" db:"password" binding:"required"`
	NickName string `json:"loginName" db:"login_name"`
}
type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
