package model

import (
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
)

type SysParamType = int

const (
	SysParamString = iota
	SysParamText
	SysParamJson
	SysParamNumber
	SysParamBoolean
)

type SysParam struct {
	ParamId   string      `json:"paramId" db:"param_id"`
	ParamKey  string      `json:"paramKey" db:"param_key"`
	ParamName string      `json:"paramName" db:"param_name"`
	CreatedAt *types.Time `json:"createdAt" db:"created_at"`
	CreatedBy string      `json:"createdBy" db:"created_by"`
	UpdatedAt *types.Time `json:"updatedAt" db:"updated_at"`
	UpdatedBy string      `json:"updatedBy" db:"updated_by"`
	Items     []ParamItem `json:"items"`
}

type ParamItem struct {
	ItemId      string       `json:"itemId" db:"item_id"`
	ParamId     string       `json:"paramId" db:"param_id"`
	ValueType   SysParamType `json:"valueType" db:"value_type"`
	ItemName    string       `json:"itemName" db:"item_name"`
	ItemKey     string       `json:"itemKey" db:"item_key"`
	ItemValue   string       `json:"itemValue" db:"item_value"`
	Description string       `json:"description" db:"description"`
	ParentId    string       `json:"parentId" db:"parent_id"`
	CreatedAt   *types.Time  `json:"createdAt" db:"created_at"`
	CreatedBy   string       `json:"createdBy" db:"created_by"`
	UpdatedAt   *types.Time  `json:"updatedAt" db:"updated_at"`
	UpdatedBy   string       `json:"updatedBy" db:"updated_by"`
	Children    []ParamItem  `json:"children"`
}

type SysParamQueryForm struct {
	web.Pageable[SysParam]
	ParamName      string                `form:"paramName"`
	ParamKey       string                `form:"paramKey"`
	ParamType      *SysParamType         `form:"paramType"`
	CreatedAtRange *types.TimeRangeArray `form:"createdAtRange"`
}
