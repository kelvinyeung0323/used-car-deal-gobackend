package model

import (
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
)

type ResType int

const (
	ResTypeCatalog ResType = iota
	ResTypeMenu
	RestTypeButton
)

type Res struct {
	ResId     string      `json:"resId" db:"res_id"`
	ParentId  string      `json:"parentId" db:"parent_id"`
	ResName   string      `json:"resName" db:"res_name"`
	ResType   ResType     `json:"resType" db:"res_type"`
	Icon      string      `json:"icon" db:"icon"`
	Sort      int64       `json:"sort" db:"sort"`
	Url       string      `json:"url" db:"url"`
	Enabled   bool        `json:"enabled" db:"enabled"`
	Visible   bool        `json:"visible" db:"visible"`
	Component string      `json:"component" db:"component"`
	AuthCode  string      `json:"authCode" db:"auth_code"`
	CreatedAt *types.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *types.Time `json:"updatedAt" db:"updated_at"`
	Children  []Res       `json:"children,omitempty"`
}

type ResQueryForm struct {
	web.Sortable[Res]
	ResType  *ResType `form:"resType"`
	ResName  string   `form:"resName"`
	ParentId string   `form:"parentId"`
	AuthCode string   `form:"authCode"`
	UserId   string   `form:"userId"`
	Enabled  *bool    `form:"enabled"`
}
