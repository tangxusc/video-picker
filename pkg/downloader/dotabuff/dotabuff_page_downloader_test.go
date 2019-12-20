package dotabuff

import (
	"context"
	"testing"
)

func TestDownload(t *testing.T) {
	downloader := NewDotaBuffPageDownloader(1)
	downloader.Download(context.TODO(), nil)
}
