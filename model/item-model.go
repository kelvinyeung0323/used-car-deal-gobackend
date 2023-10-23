package model

import (
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
)

type Item struct {
	ItemId               string      `json:"itemId" db:"item_id"` //商品编号
	Title                string      `json:"title" db:"title"`
	SubTitle             string      `json:"subTitle" db:"sub_title"`
	VehicleModel         string      `json:"vehicleModel" db:"vehicle_model"`
	UsageNature          string      `json:"usageNature" db:"usage_nature"`
	EmissionStd          string      `json:"emissionStd" db:"emission_std"`
	Displacement         string      `json:"displacement" db:"displacement"`
	Gearbox              string      `json:"gearbox" db:"gearbox"`
	BodyColor            string      `json:"bodyColor" db:"body_color"`
	BelongingPlace       string      `json:"belongingPlace" db:"belonging_place"`
	Mileage              string      `json:"mileage" db:"mileage"`
	RegistrationDate     string      `json:"registrationDate" db:"registration_date"`
	AnnualCheckDate      string      `json:"annualCheckDate" db:"annual_check_date"`
	InsuranceExpiredDate string      `json:"insuranceExpiredDate" db:"insurance_expired_date"`
	UsageTaxValidDate    string      `json:"usageTaxValidDate" db:"usage_tax_valid_date"`
	OnSale               bool        `json:"onSale" db:"on_sale"`
	AutoOnSale           bool        `json:"autoOnSale" db:"auto_on_sale"`
	StartDate            *types.Time `json:"startDate" db:"start_date"`
	EndDate              *types.Time `json:"endDate" db:"end_date"`
	CreatedAt            *types.Time `json:"createdAt" db:"created_at"`
	CreatedBy            string      `json:"createdBy" db:"created_by"`
	UpdatedAt            *types.Time `json:"updatedAt" db:"updated_at"`
	UpdatedBy            string      `json:"updatedBy" db:"updated_by"`
	Medias               []ItemMedia `json:"medias"`
	Remark               string      `json:"remark" db:"remark"`
	Deleted              bool        `json:"-" db:"deleted"`
}

type MediaType int

const (
	MediaImage MediaType = iota
	MediaVideos
	MediaFile
)

type ItemMedia struct {
	MediaId   string      `json:"mediaId" db:"media_id"`
	ItemId    string      `json:"itemId" db:"item_id"`
	Sort      string      `json:"sort" db:"sort"`
	MediaType MediaType   `json:"mediaType" db:"media_type"`
	Thumbnail string      `json:"thumbnail" db:"thumbnail"`
	Location  string      `json:"location" db:"location"`
	CreatedAt *types.Time `json:"createdAt" db:"created_at"`
	Deleted   bool        `json:"-" db:"deleted"`
}

type CheckReport struct {
}

type ItemQueryForm struct {
	web.Pageable[Item]
	ItemId           string                `form:"itemId"`
	Title            string                `form:"title"`
	SubTitle         string                `form:"subTitle"`
	UsageNature      string                `form:"usageNature"`
	EmissionStd      string                `form:"emissionStd"`
	Displacement     string                `form:"displacement"`
	Gearbox          string                `form:"gearbox"`
	BodyColor        string                `form:"bodyColor"`
	BelongingPlace   string                `form:"belongingPlace"`
	Mileage          string                `form:"mileage"`
	StartAt          *types.Time           `form:"startAt"`
	EndAt            *types.Time           `form:"endAt"`
	CreatedTimeRange *types.TimeRangeArray `form:"createdTimeRange"`
}
