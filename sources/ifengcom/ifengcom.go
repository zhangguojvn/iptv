package ifengcom

import (
	"bytes"
	"html/template"
	"iptv/common/config"
	"iptv/common/epg"
	"iptv/common/playlist"
	"iptv/sources"
	"regexp"
	"strings"

	colly "github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

type IfengCom struct{}

func (i IfengCom) Download(c config.ChannelSourceConfig) (result []playlist.Channel, epgResult []epg.EPG, err error) {
	log.Info("IfengCom Download start...")
	var epgxml epg.EPGXML
	var epgRAWs []epg.EPG
	logoTemplate, err := template.New("logo").Parse(c.Logo.Url)
	if err != nil {
		return
	}
	epgxml, err = epg.DownloadEPG(c.EPG.Url)
	if err != nil {
		log.Error(err)
		return
	}
	epgRAWs, err = epg.EPGXMLTOEPG(epgxml)
	if err != nil {
		log.Error(err)
		return
	}
	epgResult = epg.Rename(epgRAWs, c.EPG.Config["rename"])
	result = i.Crawler(c)
	log.Info("Find ", len(result), " channels.")
	for i, r := range result {
		var maplogo string
		var logoBuffer bytes.Buffer
		if c.Config["rename"] != nil && c.Config["rename"][r.Name] != nil {
			result[i].Name = c.Config["rename"][r.Name].(string)
		}
		if c.Logo.Config["link"] != nil && c.Logo.Config["link"][result[i].Name] != nil {
			maplogo = c.Logo.Config["link"][result[i].Name].(string)
		} else {

			err := logoTemplate.Execute(&logoBuffer, struct{ Name string }{Name: result[i].Name})
			if err == nil {
				maplogo = logoBuffer.String()
			}
		}
		for _, info := range epgResult {
			if info.Name == result[i].Name {
				result[i].Tvg.Id = info.Channel.ID
				result[i].Tvg.Name = info.Name
				result[i].Logo = maplogo
			}
		}
	}
	return

}
func (i IfengCom) Field(c config.ChannelSourceConfig, RAWChannels []playlist.Channel, epgInfos []epg.EPG) (groups []playlist.PlayList, epgField []epg.EPG, err error) {

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

func (i IfengCom) Crawler(con config.ChannelSourceConfig) (result []playlist.Channel) {
	var jsregex, m3uregex string
	if con.Config["crawler"] != nil && con.Config["crawler"]["jsregex"] != nil {
		jsregex = con.Config["crawler"]["jsregex"].(string)
	} else {
		jsregex = "^https://x2\\.ifengimg\\.com/fe/shank/channel/vendors\\..*\\.js"
	}
	if con.Config["crawler"] != nil && con.Config["crawler"]["m3uregex"] != nil {
		m3uregex = con.Config["crawler"]["m3uregex"].(string)
	} else {
		m3uregex = "^.*:\"http://.*\\.m3u8\"$"
	}
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		if ok, err := regexp.Match(jsregex, []byte(e.Attr("src"))); ok && err == nil {
			log.Info("Find the js :", e.Attr("src"))
			e.Request.Visit(e.Attr("src"))
		}

	})

	c.OnResponse(func(r *colly.Response) {
		log.Info("Response From : ", r.Request.URL.String())
		if ok, err := regexp.Match(jsregex, []byte(r.Request.URL.String())); ok && err == nil {
			result = i.GetM3uFromJS(string(r.Body), m3uregex)
		}
	})
	c.Visit(con.Url)
	return
}

func (IfengCom) GetM3uFromJS(js string, m3uregex string) (result []playlist.Channel) {
	var m3uCache map[string]int = make(map[string]int)
	var m3umap map[string][]string = make(map[string][]string)
	js = strings.Replace(strings.Replace(js, "{", ",", -1), "}", ",", -1)
	jsCommands := strings.Split(js, ",")
	for _, jsCommand := range jsCommands {
		if ok, err := regexp.Match(m3uregex, []byte(jsCommand)); ok && err == nil {
			formatedJsCommand := strings.Replace(strings.Replace(jsCommand, "\"", "", -1), ":", " ", 1)

			m3uList := strings.Split(formatedJsCommand, " ")
			if len(m3uList) != 2 {
				return
			}
			log.Info("Find ", m3uList[0], "'s m3u8 :", m3uList[1])
			if m3uCache[m3uList[1]] == 0 {
				m3uCache[m3uList[1]]++
				m3umap[m3uList[0]] = append(m3umap[m3uList[0]], m3uList[1])
			}
		}
	}
	for name, urls := range m3umap {
		result = append(result, playlist.Channel{
			Name: name,
			Urls: urls,
		})
	}
	return

}

func init() {
	sources.SourceRegister("ifengcom", IfengCom{})
}
