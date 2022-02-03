package sources

import (
	"iptv/common/config"
	"iptv/common/epg"
	"iptv/common/playlist"
)

func SourceRegister(name string, c SourceDownloader) {
	mapSourceDownloader[name] = c
}
func GetSource(name string) SourceDownloader {
	return mapSourceDownloader[name]
}

type SourceDownloader interface {
	Download(config.ChannelSourceConfig) (result []playlist.Channel, epgResult []epg.EPG, err error)         //生成m3u8列表
	Field(config.ChannelSourceConfig, []playlist.Channel, []epg.EPG) ([]playlist.PlayList, []epg.EPG, error) //生成组
}

var mapSourceDownloader map[string]SourceDownloader = make(map[string]SourceDownloader)
