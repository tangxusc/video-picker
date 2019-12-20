package dotabuff

import (
	"context"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/video-picker/pkg/config"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type DotaBuffDownloader struct {
	timeout int
	baseUrl string
}

func NewDotaBuffDownloader() *DotaBuffDownloader {
	return &DotaBuffDownloader{
		timeout: config.Instance.Downloader.TimeOut,
		baseUrl: "https://www.dotabuff.com/clips/%s",
	}
}

func (h *DotaBuffDownloader) Download(ctx context.Context, values map[string]interface{}) error {
	target := values[`target`].(string)
	doc, e := htmlquery.LoadURL(fmt.Sprintf(h.baseUrl, target))
	if e != nil {
		return e
	}
	one := htmlquery.FindOne(doc, "//meta[@name='twitter:player']")
	var downloadUrl string
	for _, attribute := range one.Attr {
		if `content` == attribute.Key {
			downloadUrl = attribute.Val
			break
		}
	}
	if downloadUrl == `` {
		return fmt.Errorf("未找到视频地址")
	}
	downloadUrl = strings.ReplaceAll(downloadUrl, `https://www.dotabuff.com/clips/`, `https://clip-media.dotabuff.com/`)
	downloadUrl = strings.ReplaceAll(downloadUrl, `/embed`, `.webm`)

	outputFileName, e := getOutputFileName(downloadUrl)
	if e != nil {
		return e
	}
	return download(ctx, downloadUrl, outputFileName, values)
}

func getOutputFileName(target string) (string, error) {
	dir := filepath.Join(config.Instance.Downloader.OutPath, `dotabuff`)
	_, e := os.Open(dir)
	if e != nil {
		if os.IsNotExist(e) {
			e = os.MkdirAll(dir, os.ModePerm)
			if e != nil {
				return ``, e
			}
		} else {
			return ``, e
		}
	}

	filename := fmt.Sprintf("%s.mp4", filepath.Join(dir, time.Now().Format("2006-01-02-15-04-05")))
	return filename, nil
}

func download(ctx context.Context, target string, out string, values map[string]interface{}) error {
	reader, writer := io.Pipe()
	ints := make(chan int)
	go func() {
		defer writer.Close()
		for {
			select {
			case <-ctx.Done():
				_, e := writer.Write([]byte(`q`))
				if e != nil {
					logrus.Errorf(`发送停止命令出现错误,详情:%v`, e)
				} else {
					logrus.Infof(`已发送停止命令`)
				}
				_, _ = writer.Write([]byte("\n"))
			case <-ints:
				return
			default:
				_, _ = writer.Write([]byte("\n"))
				time.Sleep(time.Second * 3)
			}
		}
	}()

	c := fmt.Sprintf("ffmpeg -y -hide_banner -i %s -strict -2 -c:v copy -c:a copy %s", target, out)
	cmd := exec.Command("sh", "-c", c)
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logrus.Infof(`开始下载:%v,输出文件:%v`, target, out)
	e := cmd.Run()
	logrus.Infof(`%v 下载完成`, target)
	ints <- 0

	values[`filepath`] = out
	return e
}
