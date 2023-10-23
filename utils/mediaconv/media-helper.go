package mediaconv

import (
	"github.com/u2takey/ffmpeg-go"
)

func GenerateThumbnail(videoPath string, thumbnailPath string) error {
	// Create a new ffmpeg command.
	err := ffmpeg_go.Input(videoPath, ffmpeg_go.KwArgs{"ss": 1}).Output(thumbnailPath, ffmpeg_go.KwArgs{"t": 1}).Run()
	return err
}
