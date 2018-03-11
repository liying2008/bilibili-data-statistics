package data

const (
	ID         = "id"
	AID        = "aid"
	VIEW       = "view"
	DANMAKU    = "danmaku"
	REPLY      = "reply"
	FAVORITE   = "favorite"
	COIN       = "coin"
	SHARE      = "share"
	NOW_RANK   = "now_rank"
	HIS_RANK   = "his_rank"
	LIKE       = "like"
	NO_REPRINT = "no_reprint"
	COPYRIGHT  = "copyright"
)

type Video struct {
	Code    int    `json:"code"`
	*Data          `json:"data"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}

type Data struct {
	Aid       uint64 `json:"aid"`
	View      int   `json:"view"`
	Danmaku   int   `json:"danmaku"`
	Reply     int   `json:"reply"`
	Favorite  int   `json:"favorite"`
	Coin      int   `json:"coin"`
	Share     int   `json:"share"`
	NowRank   int   `json:"now_rank"`
	HisRank   int   `json:"his_rank"`
	Like      int   `json:"like"`
	NoReprint int   `json:"no_reprint"`
	Copyright int   `json:"copyright"`
}
