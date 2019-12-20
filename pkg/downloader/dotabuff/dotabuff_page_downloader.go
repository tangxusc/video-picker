package dotabuff

import (
	"context"
	"fmt"
	"github.com/antchfx/htmlquery"
)

type DotaBuffPageDownloader struct {
	baseUrl   string
	pageStart uint
	xpath     string
}

func NewDotaBuffPageDownloader(startPage uint) *DotaBuffPageDownloader {
	return &DotaBuffPageDownloader{
		baseUrl:   "https://www.dotabuff.com/clips/explore?order=recent&page=%v",
		pageStart: startPage,
		xpath:     `/html/body/div[1]/div[2]/div/div/div/section/article/div[1]/div[%v]/div[1]/div[1]/a`,
	}
}

func (h *DotaBuffPageDownloader) Download(ctx context.Context, values map[string]interface{}) error {
	//target := values[`target`].(string)
	//doc, e := htmlquery.LoadURL(fmt.Sprintf(h.baseUrl, target))
	//if e != nil {
	//	return e
	//}

	doc, e := htmlquery.LoadDoc(`/home/tangxu/openProject/video-picker/video/dotabuff/source.html`)
	if e != nil {
		return e
	}
	divs := htmlquery.Find(doc, `/html/body/div[1]/div[2]/div/div/div/section/article/div[1]`)
	for _, div := range divs {
		find := htmlquery.Find(div, `a`)
		fmt.Println(find)
	}
	return nil
}
