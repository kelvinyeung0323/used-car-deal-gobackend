package model

type UploadedFile struct {
	Id           string `json:"id""`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	DeletedUrl   string `json:"deletedUrl"`
	DeleteType   string `json:"deleteType"`
}
