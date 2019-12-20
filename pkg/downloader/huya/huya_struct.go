package huya

type GameStreamInfo struct {
	SCdnType      string `json:"sCdnType"`
	IIsMaster     int    `json:"iIsMaster"`
	LChannelId    int    `json:"lChannelId"`
	LSubChannelId int    `json:"lSubChannelId"`
	LPresenterUid int    `json:"lPresenterUid"`
	SStreamName   string `json:"sStreamName"`
	SHlsUrl       string `json:"sHlsUrl"`
	SHlsUrlSuffix string `json:"sHlsUrlSuffix"`
	SHlsAntiCode  string `json:"sHlsAntiCode"`
}

type GameLiveInfo struct {
	Nick string `json:"nick"`
}

type StreamInfo struct {
	GameLiveInfo       *GameLiveInfo    `json:"gameLiveInfo"`
	GameStreamInfoList []GameStreamInfo `json:"gameStreamInfoList"`
}

type MultiStreamInfo struct {
	SDisplayName string `json:"sDisplayName"`
	IBitRate     int    `json:"iBitRate"`
}

type Stream struct {
	Status           int               `json:"status"`
	Msg              string            `json:"msg"`
	Data             []StreamInfo      `json:"data"`
	VMultiStreamInfo []MultiStreamInfo `json:"vMultiStreamInfo"`
}

type HyPlayerConfig struct {
	Html5     int     `json:"html5"`
	WEBYYHOST string  `json:"WEBYYHOST"`
	WEBYYSWF  string  `json:"WEBYYSWF"`
	WEBYYFROM string  `json:"WEBYYFROM"`
	Vappid    int     `json:"vappid"`
	Stream    *Stream `json:"stream"`
}
