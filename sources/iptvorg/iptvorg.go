package iptvorg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"iptv/common/config"
	"iptv/common/epg"
	"iptv/common/playlist"
	"iptv/sources"
	"net/http"
	"regexp"
	"text/template"
)

type IPTVorg struct {
}

func (i IPTVorg) Download(c config.ChannelSourceConfig) (result []playlist.Channel, epgResult []epg.EPG, err error) {
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
	var ChannelsMap map[string][]IPTVorgJSON = make(map[string][]IPTVorgJSON)
	logoTemplate, err := template.New("logo").Parse(c.Logo.Url)
	if err != nil {
		return
	}
	DownloadRawChannels, err := i.DownloadChannel(c.Url)
	if err != nil {
		return
	}
	for _, RawChannel := range DownloadRawChannels {
		ChannelsMap[RawChannel.Name] = append(ChannelsMap[RawChannel.Name], RawChannel)
	}
	for name, RawChannels := range ChannelsMap {
		if len(RawChannels) == 0 {
			continue
		}
		var mapname, maplogo string
		var urls = make([]string, len(RawChannels))
		var logoBuffer bytes.Buffer
		if c.Config["rename"] != nil && c.Config["rename"][name] != nil {
			mapname = c.Config["rename"][name].(string)
		} else {
			mapname = name
		}
		if c.Logo.Config["link"] != nil && c.Logo.Config["link"][mapname] != nil {
			maplogo = c.Logo.Config["link"][mapname].(string)
		} else {

			err := logoTemplate.Execute(&logoBuffer, struct{ Name string }{Name: mapname})
			if err == nil {
				maplogo = logoBuffer.String()
			}
		}
		for i, RawChannel := range RawChannels {
			urls[i] = RawChannel.Url
		}
		for _, epgInfo := range epgResult {
			if mapname == epgInfo.Name {
				RawChannels[0].Tvg.Id = epgInfo.Channel.ID
				RawChannels[0].Tvg.Name = epgInfo.Name
			}
		}
		result = append(result, playlist.Channel{
			Name:      mapname,
			Logo:      maplogo,
			Urls:      urls,
			Category:  RawChannels[0].Category,
			Languages: RawChannels[0].Languages,
			Countries: RawChannels[0].Countries,
			Tvg:       RawChannels[0].Tvg,
		})
	}
	return
}

func (i IPTVorg) Field(c config.ChannelSourceConfig, RAWChannels []playlist.Channel, epgInfos []epg.EPG) (groups []playlist.PlayList, epgField []epg.EPG, err error) {

	for _, configGroup := range c.Groups {
		var group = playlist.PlayList{
			Name: configGroup.Name,
		}
		for _, channel := range RAWChannels {
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
			var haveKeyWordUrls []string
			var ipUrls []string
			var domainUrls []string
			for _, url := range channel.Urls {
				if !c.Test.Enable || playlist.TestM3u8(url, c.Test.Timeout) {
					if configGroup.Field.Sources.UrlKeyWord != "" {
						if ok, err := regexp.Match(configGroup.Field.Sources.UrlKeyWord, []byte(url)); ok && err != nil {
							haveKeyWordUrls = append(haveKeyWordUrls, url)
							continue
						}
					}
					if ok, err := regexp.Match("([0-9]{0,3}\\.){3}[0-9]{0,3}(:|/)", []byte(url)); !ok || err != nil {
						domainUrls = append(domainUrls, url)
					} else {
						ipUrls = append(ipUrls, url)
					}

				}
			}
			if configGroup.Field.Sources.UrlType == "ip" {
				if len(haveKeyWordUrls) == 0 && len(ipUrls) == 0 {
					if c.Config["default"] != nil && c.Config["default"][channel.Name] != nil {
						if !c.Test.Enable || playlist.TestM3u8(c.Config["default"][channel.Name].(string), c.Test.Timeout) {
							channel.Urls = append([]string{c.Config["default"][channel.Name].(string)}, domainUrls...)
						}
					} else {
						channel.Urls = domainUrls
					}
				} else {
					channel.Urls = append(haveKeyWordUrls, append(ipUrls, domainUrls...)...)
				}
			} else {
				if len(haveKeyWordUrls) == 0 && len(domainUrls) == 0 {
					if c.Config["default"] != nil && c.Config["default"][channel.Name] != nil {
						if playlist.TestM3u8(c.Config["default"][channel.Name].(string), c.Test.Timeout) {
							channel.Urls = append([]string{c.Config["default"][channel.Name].(string)}, ipUrls...)
						}
					} else {
						channel.Urls = ipUrls
					}
				} else {
					channel.Urls = append(haveKeyWordUrls, append(domainUrls, ipUrls...)...)
				}

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

func (IPTVorg) DownloadChannel(u string) (result []IPTVorgJSON, err error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
	req.Header.Set("Cookie", "CONSENT=YES+cb; YSC=DwKYllHNwuw")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &result)
	return
}

func init() {
	sources.SourceRegister("iptv-org", IPTVorg{})
}
