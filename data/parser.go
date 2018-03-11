package data

import "encoding/json"

func ParseVideoData(jsonStr string) *Video {
	video := &Video{}
	json.Unmarshal([]byte(jsonStr), video)
	return video
}
