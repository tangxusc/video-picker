package downloader

type Downloader interface {
	Download(target interface{})
}
