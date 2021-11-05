package main

import (
	"fmt"
	"iptv/common/must"
	"iptv/config"
	"iptv/epg"
	"iptv/iptv"
	"iptv/m3u"
	"net/http"
	"net/url"
	"os"
)

func DownGroup() {
	var m3ugorup []m3u.M3UGroup = []m3u.M3UGroup{}
	var resultxml []epg.EPG
	for _, group := range config.Config.Groups {
		mg := m3u.M3UGroup{Group: group.DisplayName}
		ex := []epg.EPG{}
		tvs := must.Must2(iptv.IPTVField(iptv.IPTV{Name: group.IPTV.Regex, UrlType: group.IPTV.UrlType, Languages: group.IPTV.Languages}, config.Config.IPTV.Source, config.Config.IPTV.Test.Enable, config.Config.IPTV.Test.Timeout)).([]iptv.IPTV)
		infos := must.Must2(epg.EPGField(group.EPG.Regex, config.Config.EPG.Source)).([]epg.EPG)
		for _, info := range infos {
			if config.Config.IPTV.Map[info.Name] != "" {
				for _, tv := range tvs {
					if tv.Name == config.Config.IPTV.Map[info.Name] {
						tv.Tvg.Id = info.Channel.ID
						tv.Tvg.Name = info.Name
						if config.Config.Logo.Source != "" {
							tv.Logo = config.Config.Logo.Source + info.Name + ".png"
						}
						mg.IPTV = append(mg.IPTV, tv)
						ex = append(ex, info)
					}
				}
			} else {
				for _, tv := range tvs {
					if tv.Name == info.Name {
						tv.Tvg.Id = info.Channel.ID
						tv.Tvg.Name = info.Name
						if config.Config.Logo.Source != "" {
							tv.Logo = config.Config.Logo.Source + info.Name + ".png"
						}
						mg.IPTV = append(mg.IPTV, tv)
						ex = append(ex, info)
					}
				}
			}
		}

		must.Must(m3u.M3UGroupToFile([]m3u.M3UGroup{mg}, fmt.Sprintf("result/%s.m3u", group.Name)))
		must.Must(epg.EPGSAVE(must.Must2(epg.EPGTOEPGXML(ex)).(epg.EPGXML), fmt.Sprintf("result/%s.xml", group.Name)))
		m3ugorup = append(m3ugorup, mg)
		resultxml = append(resultxml, ex...)
	}
	must.Must(m3u.M3UGroupToFile(m3ugorup, fmt.Sprintf("result/%s.m3u", "all")))
	must.Must(epg.EPGSAVE(must.Must2(epg.EPGTOEPGXML(resultxml)).(epg.EPGXML), fmt.Sprintf("result/%s.xml", "all")))
}
func main() {
	config.LoadConfig("config.json")
	if config.Config.Proxy != "" {
		proxyUrl := must.Must2(url.Parse(config.Config.Proxy)).(*url.URL)
		http.DefaultClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	must.Must(os.MkdirAll("result", 0777))
	DownGroup()
}
