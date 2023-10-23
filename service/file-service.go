package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"used-car-deal-gobackend/base/config"
)

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

func (s *FileService) MoveToAvatarByUrl(ctx *gin.Context, url string) (string, error) {
	srcPath := s.TransToFilePath(url)
	return s.MoveToAvatar(ctx, srcPath)
}
func (s *FileService) MoveToAvatar(ctx *gin.Context, srcPath string) (string, error) {
	dir := config.FileUploadConfig.FileStorage + "/Avatars"
	fileName := srcPath[strings.LastIndex(srcPath, "/"):]
	dst := dir + fileName
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}
	err := os.Rename(srcPath, dst)
	return dst, err
}

func (s *FileService) MoveToDirByUrl(ctx *gin.Context, url string, targetDir string) (string, error) {
	srcPath := s.TransToFilePath(url)
	return s.MoveToDir(ctx, srcPath, targetDir)
}
func (s *FileService) MoveToDir(ctx *gin.Context, srcPath string, targetDir string) (string, error) {
	dir := config.FileUploadConfig.FileStorage + targetDir
	fileName := srcPath[strings.LastIndex(srcPath, "/"):]
	dst := dir + fileName
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}
	err := os.Rename(srcPath, dst)
	return dst, err
}

func (s *FileService) MoveToImages(ctx *gin.Context, srcPath string) (string, error) {
	dir := "/Images"
	return s.MoveToDir(ctx, srcPath, dir)
}

func (s *FileService) CopyToImages(ctx *gin.Context, srcPath string) (string, error) {
	//按月划分目录
	ts := time.Now().Format("200601")
	dir := config.FileUploadConfig.FileStorage + "/Images/" + ts
	fileName := srcPath[strings.LastIndex(srcPath, "/"):]
	dst := dir + fileName
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}
	f1, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer f1.Close()
	f2, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer f2.Close()
	_, err = io.Copy(f2, f1)
	return dst, err
}
func (s *FileService) MoveToFiles(ctx *gin.Context, srcPath string) (string, error) {
	dir := "/Files"
	return s.MoveToDir(ctx, srcPath, dir)
}

func (s *FileService) MoveToVideos(ctx *gin.Context, srcPath string) (string, error) {
	dir := "/Videos"
	return s.MoveToDir(ctx, srcPath, dir)
}

func (s *FileService) MoveToItemMedia(ctx *gin.Context, itemId string, srcPath string) (string, error) {
	//按月划分目录
	ts := time.Now().Format("200601")
	dir := config.FileUploadConfig.FileStorage + "/ItemMedia/" + ts
	fileName := srcPath[strings.LastIndex(srcPath, "/"):]
	dst := dir + fileName
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}
	err := os.Rename(srcPath, dst)
	return dst, err
}

func (s *FileService) CopyToItemMedia(ctx *gin.Context, srcPath string) (string, error) {
	//按月划分目录
	ts := time.Now().Format("200601")
	dir := config.FileUploadConfig.FileStorage + "/ItemMedia/" + ts
	fileName := srcPath[strings.LastIndex(srcPath, "/"):]
	dst := dir + fileName
	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}
	f1, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer f1.Close()
	f2, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer f2.Close()
	_, err = io.Copy(f2, f1)
	return dst, err
}
func (s *FileService) MoveToOthers(ctx *gin.Context, srcPath string) (string, error) {
	dir := "/Others/"
	return s.MoveToDir(ctx, srcPath, dir)
}

func (s *FileService) MoveToRecycleBin(ctx *gin.Context, srcPath string) (string, error) {
	dir := "/RecycleBin/"
	return s.MoveToDir(ctx, srcPath, dir)
}

func (s *FileService) SaveToTmp(ctx *gin.Context, file *multipart.FileHeader) (string, string, error) {
	extName := file.Filename[strings.LastIndex(file.Filename, "."):]
	//保存存到临时文件夹;
	//文件名为年月日+uuid 以时间为划分目录2006-01-02 15:04:05
	ts := time.Now().Format("20060102")
	uid := uuid.New()
	fileName := ts + "-" + uid.String() + extName
	dst := config.FileUploadConfig.FileStorage + "/Temp/" + fileName
	err := ctx.SaveUploadedFile(file, dst)
	return fileName, dst, err
}

func (s *FileService) TransToUrl(srcPath string) string {
	relativePath := strings.TrimPrefix(srcPath, config.FileUploadConfig.FileStorage)
	return config.FileUploadConfig.Host + config.FileUploadConfig.StaticFsPath + relativePath
}
func (s *FileService) TransToFilePath(url string) string {
	filePath := strings.TrimPrefix(url, config.FileUploadConfig.Host+config.FileUploadConfig.StaticFsPath)
	return config.FileUploadConfig.FileStorage + filePath
}

func (s *FileService) RemoveFileByUrl(url string) error {
	filePath := s.TransToFilePath(url)
	if filePath == "/" || filePath == "" {
		return fmt.Errorf("路径不能为根目录或为空")
	}
	return os.Remove(filePath)
}

func (s *FileService) RemoveFileByPath(path string) error {
	//防止错误删除根目录
	if path == "/" || path == "" {
		return fmt.Errorf("路径不能为根目录或为空")
	}
	return os.Remove(path)
}
func (s *FileService) ClearTemp(ctx *gin.Context) error {
	tmpDir := config.FileUploadConfig.FileStorage + "/Temp"
	err := filepath.Walk(tmpDir, func(name string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if time.Now().Sub(info.ModTime()).Hours() > 24 {
			log.Warnf("delete temp file:%v", name)
			//err:= os.Remove(name)
			return err
		}
		return nil
	})
	log.Errorf("clear temp file error:%v", err)
	return err
}
