package model

import "used-car-deal-gobackend/base/types"

type Item struct {
	ItemId               string      `json:"itemId" db:"item_id"`
	Title                string      `json:"title" db:"title"`
	SubTitle             string      `json:"subTitle" db:"sub_title"`
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
	UpdatedAt            *types.Time `json:"updatedAt" db:"updated_at"`
}

type ItemMedia struct {
	MediaId       string      `json:"mediaId" db:"media_id"`
	ItemId        string      `json:"itemId" db:"item_id"`
	Sort          string      `json:"sort" db:"sort"`
	MediaType     int         `json:"mediaType" db:"media_type"`
	MediaLocation string      `json:"mediaLocation" db:"media_location"`
	CreatedAt     *types.Time `json:"createdAt" db:"created_at"`
}

type ItemSpec struct {
	SpecId    string      `json:"specId" db:"spec_id"`
	Year      string      `json:"year" db:"year"`
	Brand     string      `json:"brand" db:"brand"`
	Model     string      `json:"model" db:"model"`
	Params    string      `json:"params" db:"params"`
	CreatedAt *types.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *types.Time `json:"updatedAt" db:"updated_at"`
}

type CheckReport struct {
}
