package controller

import (
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
	"used-car-deal-gobackend/service"
)

type FileUploadController struct {
	fileService *service.FileService
}

func NewFileUploadController() *FileUploadController {
	return &FileUploadController{}
}

func (c *FileUploadController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Warnf("文件上传失败,%v", err)
		web.Err(web.ERROR)
	}

	//区分文件类型图片放到Images/文件放到Files

	log.Debug("upload file:", file.Filename)

	if err != nil {
		log.Warnf("文件上传失败:%v", err)
		web.Err(web.ERROR)
	}
	fileName, filePath, err2 := c.fileService.SaveToTmp(ctx, file)
	if err2 != nil {
		log.Warnf("文件保存到临时文件夹失败:%v", err2)
		web.Err(web.ERROR)
	}
	uf := &model.UploadedFile{Name: fileName, Url: c.fileService.TransToUrl(filePath), Size: file.Size}

	web.ReturnOK(ctx, uf)

}
