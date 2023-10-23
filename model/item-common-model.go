package model

import (
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
)

type Brand struct {
	BrandId   int         `json:"brandId" db:"brand_id"`
	BrandLogo string      `json:"brandLogo" db:"brand_logo"`
	BrandName string      `json:"brandName" db:"brand_name"`
	Country   string      `json:"country" db:"country"`
	Remark    string      `json:"remark" db:"remark"`
	CreatedAt *types.Time `json:"createdAt" db:"created_at"`
	CreatedBy string      `json:"createdBy" db:"created_by"`
	UpdatedAt *types.Time `json:"updatedAt" db:"updated_at"`
	UpdatedBy string      `json:"updatedBy" db:"updated_by"`
	Deleted   bool        `json:"deleted" db:"deleted"`
	Series    []Series    `json:"series" db:"-"`
}

type BrandQueryForm struct {
	web.Pageable[Brand]
	BrandName      string                `form:"brandName"`
	Countries      []string              `form:"countries"`
	CreatedAtRange *types.TimeRangeArray `form:"CreatedAtRange"`
}

type Series struct {
	Image      string      `json:"image" db:"image"`
	SeriesId   int         `json:"seriesId" db:"series_id"`
	SeriesName string      `json:"seriesName" db:"series_name"`
	BrandId    int         `json:"brandId" db:"brand_id"`
	BrandName  string      `json:"brandName"`
	Grade      string      `json:"grade" db:"grade"`
	Remark     string      `json:"remark" db:"remark"`
	CreatedAt  *types.Time `json:"createdAt" db:"created_at"`
	CreatedBy  string      `json:"createdBy" db:"created_by"`
	UpdatedAt  *types.Time `json:"updatedAt" db:"updated_at"`
	UpdatedBy  string      `json:"updatedBy" db:"updated_by"`
	Deleted    string      `json:"deleted" db:"deleted"`
}
type SeriesQueryForm struct {
	web.Pageable[Brand]
	SeriesName     string                `form:"seriesName"`
	Brands         []string              `form:"brands"`
	Grades         []string              `form:"grades"`
	CreatedAtRange *types.TimeRangeArray `form:"CreatedAtRange"`
}
