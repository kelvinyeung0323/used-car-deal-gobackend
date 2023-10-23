package model

import (
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
)

type Role struct {
	RoleId    string      `json:"roleId" db:"role_id"`
	RoleName  string      `json:"roleName" db:"role_name"`
	AuthCode  string      `json:"authCode" db:"auth_code"`
	Enabled   bool        `json:"enabled" db:"enabled"`
	Remark    string      `json:"remark" db:"remark"`
	CreatedAt *types.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *types.Time `json:"updatedAt" db:"updated_at"`
	Res       []Res       `json:"res"`
	ResIds    []string    `json:"resIds"`
}

type RoleQuerySorts struct {
	RoleId    string      `json:"roleId" db:"role_id"`
	RoleName  string      `json:"roleName" db:"role_name"`
	AuthCode  string      `json:"authCode" db:"auth_code"`
	Enabled   *bool       `json:"enabled" db:"enabled"`
	Remark    string      `json:"remark" db:"remark"`
	CreatedAt *types.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *types.Time `json:"updatedAt" db:"updated_at"`
}

type RoleQueryForm struct {
	web.Pageable[RoleQuerySorts]
	RoleName string      `form:"roleName"`
	AuthCode string      `form:"authCode"`
	Enabled  *bool       `form:"enabled"`
	StartAt  *types.Time `form:"startAt"`
	EndAt    *types.Time `form:"endAt"`
}
type RoleUserQueryForm struct {
	web.Pageable[RoleQuerySorts]
	IsBelongToRole bool        `form:"isBelongToRole"` //IsBelongToRole标识，为true表示查询属于角色的用户，false 表示查询不属于角色的用户
	RoleId         string      `form:"roleId"`
	UserName       string      `form:"userName"`
	UserId         string      `form:"userId"`
	Enabled        *bool       `form:"enabled"`
	StartAt        *types.Time `form:"startAt"`
	EndAt          *types.Time `form:"endAt"`
}

type EnabledRoleForm struct {
	RoleId  string `form:"roleId" binding:"required"`
	Enabled bool   `form:"enabled" `
}

type UsersOfRole struct {
	RoleId  string   `form:"roleId"`
	UserIds []string `form:"userIds"`
}
