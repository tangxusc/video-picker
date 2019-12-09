package downloader

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/video-picker/pkg/config"
	"github.com/tangxusc/video-picker/pkg/dispatcher"
	"github.com/tangxusc/video-picker/pkg/eventbus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const baseUrl = "https://www.huya.com/%s"

type HuyaDownloader struct {
	dispatcher  *dispatcher.Dispatcher
	bus         *eventbus.Bus
	nextChannel string
}

func (h *HuyaDownloader) Download(target string) *dispatcher.Job {
	job := dispatcher.NewJob(func(ctx context.Context) error {
		return h.downM3u8(ctx, target, config.Instance.Downloader.TimeOut)
	})
	h.dispatcher.Dispatch(job)
	return job
}

func (h *HuyaDownloader) DownloadWithTimeOut(target string, timeout int) {
	job := dispatcher.NewJob(func(ctx context.Context) error {
		return h.downM3u8(ctx, target, timeout)
	})
	h.dispatcher.Dispatch(job)
}

func NewHuyaDownloader(ctx context.Context, maxWorkers int, bus *eventbus.Bus) *HuyaDownloader {
	dis := dispatcher.NewDispatcher(ctx, maxWorkers)
	downloader := &HuyaDownloader{
		dispatcher:  dis,
		bus:         bus,
		nextChannel: "downloaded",
	}
	downloader.dispatcher.Start()
	//c := make(chan interface{})
	//bus.Subscribe("downloading", c)
	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		case event := <-c:
	//			switch event.(type) {
	//			case string:
	//				downloader.Download(event.(string))
	//			}
	//		}
	//	}
	//}()
	return downloader
}

func (h *HuyaDownloader) downM3u8(ctx context.Context, target string, timeout int) error {
	hyPlayerConfig, e := h.getHyPlayerConfig(target, ctx)
	if e != nil {
		return e
	}
	dir := filepath.Join(config.Instance.Downloader.OutPath, hyPlayerConfig.Stream.Data[0].GameLiveInfo.Nick)
	_, e = os.Open(dir)
	if e != nil {
		if os.IsNotExist(e) {
			e = os.MkdirAll(dir, os.ModePerm)
			if e != nil {
				return e
			}
		} else {
			return e
		}
	}

	var m3u8 string
	for _, v := range hyPlayerConfig.Stream.Data[0].GameStreamInfoList {
		if v.SHlsUrlSuffix == "m3u8" {
			m3u8 = fmt.Sprintf("%s/%s.%s", v.SHlsUrl, v.SStreamName, v.SHlsUrlSuffix)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			e := h.download(ctx, m3u8, dir, timeout)
			if e != nil {
				return e
			}
		}
	}
}

func (h *HuyaDownloader) getHyPlayerConfig(target string, ctx context.Context) (*HyPlayerConfig, error) {
	url := fmt.Sprintf(baseUrl, target)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequestWithContext(ctx, `GET`, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	htmlStr := string(body)
	tmp := strings.Split(htmlStr, "hyPlayerConfig =")
	if len(tmp) < 2 {
		return nil, fmt.Errorf("解析hyPlayerConfig异常")
	}

	tmp = strings.Split(tmp[1], "window.TT_LIVE_TIMING")
	if len(tmp) < 2 {
		return nil, fmt.Errorf("解析window.TT_LIVE_TIMING异常")
	}
	jsonStr := strings.Replace(tmp[0], "};", "}", 1)
	var hyPlayerConfig HyPlayerConfig
	err = json.Unmarshal([]byte(jsonStr), &hyPlayerConfig)
	if err != nil {
		return nil, err
	}
	if hyPlayerConfig.Stream == nil || len(hyPlayerConfig.Stream.Data) <= 0 {
		return nil, fmt.Errorf(`解析异常或主播还未开播`)
	}
	return &hyPlayerConfig, nil
}

func (h *HuyaDownloader) download(ctx context.Context, target string, dir string, timeout int) error {
	reader, writer := io.Pipe()
	go func() {
		ticker := time.NewTicker(time.Duration(timeout) * time.Minute)
		defer ticker.Stop()
		defer writer.Close()
		send := false
		select {
		case <-ticker.C:
			send = true
		case <-ctx.Done():
			send = true
		}
		if send {
			_, e := writer.Write([]byte(`q`))
			if e != nil {
				logrus.Errorf(`发送停止命令出现错误,详情:%v`, e)
			} else {
				logrus.Infof(`已发送停止命令`)
			}
			_, _ = writer.Write([]byte("\n"))
		}
	}()
	logrus.Infof(`开始下载:%v,输出目录:%v`, target, dir)
	filename := fmt.Sprintf("%s.mp4", filepath.Join(dir, time.Now().Format("2006-01-02-15-04-05")))
	c := fmt.Sprintf("ffmpeg -y -hide_banner -loglevel info -i %s -c:v copy -c:a copy %s", target, filename)
	cmd := exec.Command("sh", "-c", c)
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Run()
	logrus.Infof(`%v 下载完成`, target)

	h.bus.Send(h.nextChannel, filename)
	return e
}
