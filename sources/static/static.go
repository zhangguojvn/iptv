package static

import (
	"iptv/common/config"
	"iptv/common/epg"
	"iptv/common/playlist"
	"iptv/sources"
	"regexp"
)

type Static struct{}

func (Static) Download(c config.ChannelSourceConfig) (result []playlist.Channel, epgResult []epg.EPG, err error) {
	var epgxml epg.EPGXML
	var epgRAWs []epg.EPG
	epgxml, err = epg.DownloadEPG(c.EPG.Url)
	if err != nil {
		return
	}
	epgRAWs, err = epg.EPGXMLTOEPG(epgxml)
	if err != nil {
		return
	}
	epgResult = epg.Rename(epgRAWs, c.EPG.Config["rename"])
	for name, link := range c.Config["default"] {
		cha := playlist.Channel{
			Name: name,
			Urls: []string{link.(string)},
		}
		for _, epgInfo := range epgResult {
			if name == epgInfo.Name {
				cha.Tvg.Id = epgInfo.Channel.ID
				cha.Tvg.Name = epgInfo.Name
			}
		}
		result = append(result, cha)
	}
	return
}
func (Static) Field(c config.ChannelSourceConfig, RAWchannels []playlist.Channel, epgInfos []epg.EPG) (groups []playlist.PlayList, epgField []epg.EPG, err error) {

	for _, configGroup := range c.Groups {
		var group = playlist.PlayList{
			Name: configGroup.Name,
		}
		for _, channel := range RAWchannels {
			if configGroup.Field.Sources.Regex != "" {
				if ok, err := regexp.Match(configGroup.Field.Sources.Regex, []byte(channel.Name)); !ok || err != nil {
					continue
				}
			}
			if configGroup.Field.Sources.Languages != nil && !func() bool {
				if channel.Languages == nil {
					return false
				}
				for _, requestL := range configGroup.Field.Sources.Languages {
					for _, haveL := range channel.Languages {
						if haveL == requestL {
							goto A
						}
					}
					return false
				A:
				}

				return true
			}() {
				continue
			}
			if configGroup.Field.Sources.Countries != nil && !func() bool {
				if channel.Countries == nil {
					return false
				}

				for requestC := range configGroup.Field.Sources.Countries {
					for haveC := range channel.Countries {
						if haveC == requestC {
							goto B
						}
					}
					return false
				B:
				}

				return true
			}() {
				continue
			}
			var ipUrls []string
			var domainUrls []string
			for _, url := range channel.Urls {
				if !c.Test.Enable || playlist.TestM3u8(url, c.Test.Timeout) {
					if ok, err := regexp.Match("([0-9]{0,3}\\.){3}[0-9]{0,3}(:|/)", []byte(url)); !ok || err != nil {
						domainUrls = append(domainUrls, url)
					} else {
						ipUrls = append(ipUrls, url)
					}

				}
			}
			if configGroup.Field.Sources.UrlType == "ip" {
				channel.Urls = append(ipUrls, domainUrls...)
			} else {
				channel.Urls = append(domainUrls, ipUrls...)
			}
			group.Channels = append(group.Channels, channel)
		}
		group.Channels = playlist.ChannelSort(group.Channels)
		groups = append(groups, group)

		tepg, err := epg.EPGField(epgInfos, configGroup.Field.EPG.Regex)
		if err != nil {

		} else {
			epgField = append(epgField, tepg...)
		}
	}

	return
}

func init() {
	sources.SourceRegister("static", Static{})
}
