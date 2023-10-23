package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sync"
	"time"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/types"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/repo"
)

type ItemService struct {
	seqPool     int //每5分钟归零
	ticker      *time.Ticker
	mutex       sync.Mutex
	itemRepo    *repo.ItemRepo
	fileService *FileService
	txMgr       *datasource.TransactionManger
}

func NewItemService(itemRepo *repo.ItemRepo, fileService *FileService, txMgr *datasource.TransactionManger) *ItemService {
	ticker := time.NewTicker(5 * time.Minute)
	itemService := &ItemService{itemRepo: itemRepo, fileService: fileService, txMgr: txMgr, ticker: ticker}
	itemService.runTicker()
	return itemService
}

// CreateItem 创建只要生成一个商品编号
// 编号规则 PD202310091219+秒 + 3位随机数
func (s *ItemService) CreateItem(ctx *gin.Context, form *model.Item) {
	s.txMgr.BeginTx(ctx)
	//2006-01-02 15:04:05
	form.ItemId = "PD" + time.Now().Format("20060102150405") + s.getSeq()

	now := types.Time(time.Now())
	form.CreatedAt = &now
	form.UpdatedAt = &now

	for i := 0; i < len(form.Medias); i++ {
		media := form.Medias[i]
		media.Location = s.fileService.TransToFilePath(media.Location)
		//考虑:重复提交或错误后重新提交，所以这里做复制不做移动
		mediaPath, err := s.fileService.CopyToItemMedia(ctx, media.Location)
		if err != nil {
			log.Debugf("保存图片错误:%v", err)
			web.Err(web.ERROR)
		}
		media.MediaId = uuid.New().String()
		media.ItemId = form.ItemId
		media.Sort = string(rune(i))
		media.Location = mediaPath
		media.CreatedAt = &now
		s.itemRepo.CreateItemMedia(ctx, &media)

	}
	s.itemRepo.CreateItem(ctx, form)
	//保存图片
	s.txMgr.CommitTx(ctx)
}

func (s *ItemService) UpdateItem(ctx *gin.Context, form *model.Item) {
	s.txMgr.BeginTx(ctx)
	//2006-01-02 15:04:05
	if form.ItemId == "" {
		web.BizErr("id 不能为空")
	}

	itemInDb := s.itemRepo.FindItemById(ctx, form.ItemId)
	if itemInDb == nil {
		web.BizErr("商品不存在！")
	}

	now := types.Time(time.Now())
	form.UpdatedAt = &now

	//媒体处理
	mediasInDb := s.itemRepo.GetMediaOfItem(ctx, form.ItemId)
	mediaMap := map[string]*model.ItemMedia{}
	//注意，range出来的mdeia是副本
	for _, media := range mediasInDb {
		mediaMap[media.Location] = &media
	}
	var appendMedias []*model.ItemMedia
	for i := 0; i < len(form.Medias); i++ {
		media := form.Medias[i]
		mediaPath := s.fileService.TransToFilePath(media.Location)

		if m1 := mediaMap[mediaPath]; m1 == nil {
			mediaPath, err := s.fileService.CopyToItemMedia(ctx, mediaPath)
			if err != nil {
				log.Debugf("保存图片失败:%v", err)
				web.Err(web.ERROR)
			}
			media.MediaId = uuid.New().String()
			media.ItemId = form.ItemId
			media.Sort = string(rune(i))
			media.Location = mediaPath
			media.CreatedAt = &now
			appendMedias = append(appendMedias, &media)
		} else {
			//如果数据库中的媒体存包括表单中的媒体，则不做创建操作，并把并从map中删除，剩下的就是要从数据库删除的图片
			delete(mediaMap, mediaPath)
		}

	}

	//删除媒体 具体文件不作删除
	for _, v := range mediaMap {
		//不是物理删除，只是设置了删除标志
		s.itemRepo.DeleteItemMedia(ctx, v)
	}
	//创建媒体
	for _, media := range appendMedias {
		s.itemRepo.CreateItemMedia(ctx, media)
	}
	//保存商品
	s.itemRepo.UpdateItem(ctx, form)
	//保存图片
	s.txMgr.CommitTx(ctx)
}

func (s *ItemService) DeleteItems(ctx *gin.Context, itemsIds []string) {
	s.txMgr.BeginTx(ctx)
	for _, itemsId := range itemsIds {
		s.itemRepo.DeleteItem(ctx, itemsId)
		//删除相关Media
		s.itemRepo.DeleteMediaOfItem(ctx, itemsId)
	}
	s.txMgr.CommitTx(ctx)
}
func (s *ItemService) QueryItems(ctx *gin.Context, form *model.ItemQueryForm) *web.Page[model.Item] {
	page := s.itemRepo.QueryItems(ctx, form)
	for i := 0; i < len(page.Data); i++ {
		it := &page.Data[i]
		it.Medias = s.itemRepo.GetMediaOfItem(ctx, it.ItemId)
		for j := 0; j < len(it.Medias); j++ {
			md := &it.Medias[j]
			md.Location = s.fileService.TransToUrl(md.Location)
		}
	}

	return page
}
func (s *ItemService) FindItem(ctx *gin.Context, itemId string) *model.Item {
	item := s.itemRepo.FindItemById(ctx, itemId)
	item.Medias = s.itemRepo.GetMediaOfItem(ctx, itemId)
	//转换视频连接
	for i := 0; i < len(item.Medias); i++ {
		m := &item.Medias[i]
		m.Location = s.fileService.TransToUrl(m.Location)
	}
	return item
}

func (s *ItemService) getSeq() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.seqPool++
	return fmt.Sprintf("%04d", s.seqPool)
}

// 开户一个定时器
func (s *ItemService) runTicker() {
	go func() {
		for {
			<-s.ticker.C
			s.mutex.Lock()
			s.seqPool = 0
			s.mutex.Unlock()
		}
	}()
}
