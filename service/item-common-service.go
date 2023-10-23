package service

import (
	"github.com/gin-gonic/gin"
	"time"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/repo"
)

type ItemCommonService struct {
	commonRepo  *repo.ItemCommonRepo
	fileService *FileService
	txMgr       *datasource.TransactionManger
}

func NewItemCommonService(commonRepo *repo.ItemCommonRepo, fileService *FileService, txMgr *datasource.TransactionManger) *ItemCommonService {
	commonService := &ItemCommonService{commonRepo: commonRepo, fileService: fileService, txMgr: txMgr}
	return commonService
}

func (s *ItemCommonService) QueryBrand(ctx *gin.Context, form *model.BrandQueryForm) *web.Page[model.Brand] {
	page := s.commonRepo.PageQueryBrands(ctx, form)
	for i := 0; i < len(page.Data); i++ {
		b := &page.Data[i]
		b.BrandLogo = s.fileService.TransToUrl(b.BrandLogo)
	}
	return page
}

func (s *ItemCommonService) FindBrand(ctx *gin.Context, brandId int) *model.Brand {
	brand := s.commonRepo.FindBrand(ctx, brandId)
	brand.BrandLogo = s.fileService.TransToUrl(brand.BrandLogo)
	return brand
}

func (s *ItemCommonService) CreateBrand(ctx *gin.Context, form *model.Brand) {
	s.txMgr.BeginTx(ctx)
	now := types.Time(time.Now())
	if form.BrandLogo != "" {

		form.BrandLogo = s.fileService.TransToFilePath(form.BrandLogo)
		newPath, err := s.fileService.CopyToImages(ctx, form.BrandLogo)
		if err != nil {
			log.Warnf("复制图片失败：%v", err)
			web.Err(web.ERROR)
		}
		form.BrandLogo = newPath
	}
	form.CreatedAt = &now
	form.UpdatedAt = &now
	s.commonRepo.CreateBrand(ctx, form)
	//保存图片
	s.txMgr.CommitTx(ctx)
}

func (s *ItemCommonService) UpdateBrand(ctx *gin.Context, form *model.Brand) {
	s.txMgr.BeginTx(ctx)
	brand := s.commonRepo.FindBrand(ctx, form.BrandId)
	if brand == nil {
		web.BizErr("品牌不存在")
	}
	if form.BrandLogo != "" {
		//如果图片路径与之前一致则不操作
		form.BrandLogo = s.fileService.TransToFilePath(form.BrandLogo)
		if form.BrandLogo != brand.BrandLogo {
			newPath, err := s.fileService.CopyToImages(ctx, form.BrandLogo)
			if err != nil {
				log.Warnf("复制图片失败：%v", err)
				web.Err(web.ERROR)
			}
			form.BrandLogo = newPath
		}
	}
	now := types.Time(time.Now())
	form.UpdatedAt = &now
	s.commonRepo.UpdateBrand(ctx, form)
	//保存图片
	s.txMgr.CommitTx(ctx)
}

func (s *ItemCommonService) DeleteBrands(ctx *gin.Context, brandIds []int) {
	s.txMgr.BeginTx(ctx)
	for _, brandId := range brandIds {
		s.commonRepo.DeleteBrand(ctx, brandId)
	}
	s.txMgr.CommitTx(ctx)
}

func (s *ItemCommonService) GetAllBrandsWithSeries(ctx *gin.Context) []model.Brand {
	brands := s.commonRepo.QueryBrands(ctx, &model.BrandQueryForm{})
	for i := 0; i < len(brands); i++ {
		brand := &brands[i]
		brand.BrandLogo = s.fileService.TransToUrl(brand.BrandLogo)
		series := s.commonRepo.GetSeriesOfBrand(ctx, brand.BrandId)
		brand.Series = series
		for j := 0; j < len(brand.Series); j++ {
			s1 := &brand.Series[j]
			s1.Image = s.fileService.TransToUrl(s1.Image)
		}
	}
	return brands
}

//==================== 车系 ===================================

func (s *ItemCommonService) QuerySeries(ctx *gin.Context, form *model.SeriesQueryForm) *web.Page[model.Series] {
	page := s.commonRepo.QuerySeries(ctx, form)
	for i := 0; i < len(page.Data); i++ {
		b := &page.Data[i]
		b.Image = s.fileService.TransToUrl(b.Image)
	}
	return page
}

func (s *ItemCommonService) FindSeries(ctx *gin.Context, seriesId int) *model.Series {
	series := s.commonRepo.FindSeries(ctx, seriesId)
	series.Image = s.fileService.TransToUrl(series.Image)
	return series
}

func (s *ItemCommonService) CreateSeries(ctx *gin.Context, form *model.Series) {
	s.txMgr.BeginTx(ctx)

	if form.Image != "" {

		form.Image = s.fileService.TransToFilePath(form.Image)
		newPath, err := s.fileService.CopyToImages(ctx, form.Image)
		if err != nil {
			log.Warnf("复制图片失败：%v", err)
			web.Err(web.ERROR)
		}
		form.Image = newPath
	}
	now := types.Time(time.Now())
	form.CreatedAt = &now
	form.UpdatedAt = &now
	s.commonRepo.CreateSeries(ctx, form)
	//保存图片
	s.txMgr.CommitTx(ctx)
}

func (s *ItemCommonService) UpdateSeries(ctx *gin.Context, form *model.Series) {
	s.txMgr.BeginTx(ctx)
	series := s.commonRepo.FindSeries(ctx, form.SeriesId)
	if series == nil {
		web.BizErr("车系不存在")
	}
	if form.Image != "" {
		//如果图片路径与之前一致则不操作
		form.Image = s.fileService.TransToFilePath(form.Image)
		if form.Image != series.Image {
			newPath, err := s.fileService.CopyToImages(ctx, form.Image)
			if err != nil {
				log.Warnf("复制图片失败：%v", err)
				web.Err(web.ERROR)
			}
			form.Image = newPath
		}
	}
	now := types.Time(time.Now())
	form.UpdatedAt = &now
	s.commonRepo.UpdateSeries(ctx, form)
	//保存图片
	s.txMgr.CommitTx(ctx)
}

func (s *ItemCommonService) DeleteSeries(ctx *gin.Context, seriesIds []int) {
	s.txMgr.BeginTx(ctx)
	for _, seriesId := range seriesIds {
		s.commonRepo.DeleteSeries(ctx, seriesId)
	}
	s.txMgr.CommitTx(ctx)
}
