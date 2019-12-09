package downloader

import "github.com/tangxusc/video-picker/pkg/dispatcher"

type Downloader interface {
	Download(target string) *dispatcher.Job
}
