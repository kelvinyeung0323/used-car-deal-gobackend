package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
	"time"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/repo"
)

type ItemSpecService struct {
	specRepo    *repo.ItemSpecRepo
	fileService *FileService
	txMgr       *datasource.TransactionManger
}

func NewItemSpecService(specRepo *repo.ItemSpecRepo, fileService *FileService, txMgr *datasource.TransactionManger) *ItemSpecService {
	itemService := &ItemSpecService{specRepo: specRepo, fileService: fileService, txMgr: txMgr}
	return itemService
}

// CreateSpec 创建规格
func (s *ItemSpecService) CreateSpec(ctx *gin.Context, form *model.ItemSpec) {
	s.txMgr.BeginTx(ctx)
	//2006-01-02 15:04:05
	form.SpecId = strings.ReplaceAll(uuid.New().String(), "-", "")
	now := types.Time(time.Now())
	form.CreatedAt = &now
	form.UpdatedAt = &now
	s.specRepo.CreateSpec(ctx, form)
	//保存图片
	s.txMgr.CommitTx(ctx)
}

func (s *ItemSpecService) UpdateSpec(ctx *gin.Context, form *model.ItemSpec) {
	s.txMgr.BeginTx(ctx)
	//2006-01-02 15:04:05
	if form.SpecId == "" {
		web.BizErr("id 不能为空")
	}

	itemInDb := s.specRepo.FindSpecById(ctx, form.SpecId)
	if itemInDb == nil {
		web.BizErr("规格不存在不存在！")
	}
	now := types.Time(time.Now())
	form.UpdatedAt = &now
	s.specRepo.UpdateSpec(ctx, form)
	s.txMgr.CommitTx(ctx)
}

func (s *ItemSpecService) DeleteSpecs(ctx *gin.Context, specIds []string) {
	s.txMgr.BeginTx(ctx)
	for _, specId := range specIds {
		s.specRepo.DeleteSpec(ctx, specId)
		//删除相关Media
		s.specRepo.DeleteMediaOfSpec(ctx, specId)
		s.specRepo.DeleteColorOfSpec(ctx, specId)
	}
	s.txMgr.CommitTx(ctx)
}
func (s *ItemSpecService) QuerySpecs(ctx *gin.Context, form *model.ItemSpecQueryForm) *web.Page[model.ItemSpec] {
	page := s.specRepo.QuerySpecs(ctx, form)
	return page
}
func (s *ItemSpecService) FindSpec(ctx *gin.Context, itemId string) *model.ItemSpec {
	spec := s.specRepo.FindSpecById(ctx, itemId)
	return spec
}

func (s *ItemSpecService) GetColorsOfSpec(ctx *gin.Context, specId string) []model.ItemSpecColor {
	colors := s.specRepo.GetColorOfSpec(ctx, specId)
	for i := 0; i < len(colors); i++ {
		color := &colors[i]
		color.Medias = s.specRepo.GetMediaOfColor(ctx, color.ColorId)
		//转换图片路径
		for j := 0; j < len(color.Medias); j++ {
			m := &color.Medias[j]
			m.Location = s.fileService.TransToUrl(m.Location)
		}

	}
	return colors
}

func (s *ItemSpecService) CreateSpecColor(ctx *gin.Context, form *model.ItemSpecColor) {
	s.txMgr.BeginTx(ctx)
	if form.SpecId == "" {
		web.BizErr("规格ID不能为空")
	}
	specInDb := s.specRepo.FindSpecById(ctx, form.SpecId)
	if specInDb == nil {
		web.BizErr("规格不存在")
	}

	form.ColorId = strings.ReplaceAll(uuid.New().String(), "-", "")
	now := types.Time(time.Now())
	form.CreatedAt = &now
	form.CreatedBy = ""
	s.specRepo.CreateSpecColor(ctx, form)
	s.txMgr.CommitTx(ctx)
}

func (s *ItemSpecService) DeleteSpecColor(ctx *gin.Context, colorIds []string) {
	s.txMgr.BeginTx(ctx)
	for _, colorId := range colorIds {
		m := s.specRepo.GetMediaOfColor(ctx, colorId)
		s.specRepo.DeleteSpecColor(ctx, colorId)
		for _, m1 := range m {
			s.specRepo.DeleteSpecMedia(ctx, m1.MediaId)
		}
	}
	s.txMgr.CommitTx(ctx)
}

func (s *ItemSpecService) FindSpecColor(ctx *gin.Context, colorId string) *model.ItemSpecColor {
	color := s.specRepo.FindSpecColorById(ctx, colorId)
	if color != nil {
		color.Medias = s.specRepo.GetMediaOfColor(ctx, colorId)
	}
	return color
}

func (s *ItemSpecService) UpdateSpecColor(ctx *gin.Context, form *model.ItemSpecColor) {
	s.txMgr.BeginTx(ctx)
	if form.SpecId == "" || form.ColorId == "" {
		web.BizErr("规格ID和颜色ID不能为空")
	}
	spec := s.specRepo.FindSpecById(ctx, form.SpecId)
	if spec == nil {
		web.BizErr("找不到相应的规格信息")
	}
	colorInDb := s.specRepo.FindSpecColorById(ctx, form.ColorId)
	if colorInDb == nil {
		web.BizErr("找不到相应的颜色信息")
	}
	s.specRepo.UpdateSpecColor(ctx, form)

	s.txMgr.CommitTx(ctx)
}

func (s *ItemSpecService) UploadSpecMedia(ctx *gin.Context, media *model.ItemSpecMedia) {
	s.txMgr.BeginTx(ctx)
	if media.SpecId == "" {
		web.BizErr("规格ID不能为空!")
	}
	if media.ColorId == "" {
		web.BizErr("颜色ID不能为空")
	}
	if media.Location == "" {
		web.BizErr("文件路径不能为空")
	}
	media.MediaId = strings.ReplaceAll(uuid.New().String(), "-", "")
	s.specRepo.CreateSpecMedia(ctx, media)
	s.txMgr.CommitTx(ctx)
}

func (s *ItemSpecService) DeleteSpecMedia(ctx *gin.Context, mediaIds []string) {
	s.txMgr.BeginTx(ctx)
	for _, mediaId := range mediaIds {
		s.specRepo.DeleteSpecMedia(ctx, mediaId)
	}
	s.txMgr.CommitTx(ctx)
}
