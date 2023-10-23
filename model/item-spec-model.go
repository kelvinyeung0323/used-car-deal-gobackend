package model

import (
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
)

type ItemSpec struct {
	SpecId     string          `json:"specId" db:"spec_id"`
	Year       string          `json:"year" db:"year"`
	BrandId    int             `json:"brandId" db:"brand_id"`
	BrandName  string          `json:"brandName" db:"brand_name"`
	SeriesId   int             `json:"seriesId" db:"series_id"`
	SeriesName string          `json:"seriesName" db:"series_name"`
	Model      string          `json:"model" db:"model"`
	Params     []ItemSpecParam `json:"params" db:"-"` //在数据库里以json形式保存
	ParamsStr  string          `json:"-" db:"params"`
	CreatedAt  *types.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt  *types.Time     `json:"updatedAt" db:"updated_at"`
}
type ItemSpecParam struct {
	Name     string          `json:"name"`
	Value    string          `json:"value"`
	Remark   string          `json:"remark"`
	Children []ItemSpecParam `json:"children"`
}

type ItemSpecColor struct {
	ColorId   string          `json:"colorId" db:"color_id"`
	SpecId    string          `json:"specId" db:"spec_id"`
	Value     string          `json:"value" db:"value"`
	Name      string          `json:"name" db:"name"`
	Medias    []ItemSpecMedia `json:"medias"`
	CreatedAt *types.Time     `json:"createdAt" db:"created_at"`
	CreatedBy string          `json:"createdBy" db:"created_by"`
}
type ItemSpecMediaCatalog = int

const (
	mediaCatalogExterior ItemSpecMediaCatalog = iota //外观
	mediaCatalogInterior                             //内饰
	mediaCatalogSpace                                //空间)
)

type ItemSpecMedia struct {
	MediaId   string               `json:"mediaId" db:"media_id"`
	SpecId    string               `json:"specId" db:"spec_id"`
	Sort      int                  `json:"sort" db:"sort"`
	Thumbnail string               `json:"thumbnail" db:"thumbnail"`
	MediaType MediaType            `json:"mediaType" db:"media_type"`
	Catalog   ItemSpecMediaCatalog `json:"catalog" db:"catalog"`
	Location  string               `json:"location" db:"location"`
	ColorId   string               `json:"colorId" db:"color_id"`
	CreatedAt *types.Time          `json:"createdAt" db:"created_at"`
	CreatedBy *types.Time          `json:"createdBy" db:"created_by"`
}

type ItemSpecQueryForm struct {
	web.Pageable[Item]
	Year           string                `form:"year"`
	BrandIds       []string              `form:"brandIds"`
	SeriesIds      []string              `form:"seriesIds"`
	Model          string                `form:"model"`
	CreatedAtRange *types.TimeRangeArray `form:"createdAtRange"`
}

type ItemSpecColorForm struct {
	SpecId string `json:"specId"`
	Colors []ItemSpecColor
}
