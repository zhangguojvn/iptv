package main

import (
	"fmt"
	"iptv/common/config"
	"iptv/common/epg"
	"iptv/common/must"
	"iptv/common/playlist"
	"iptv/sources"
	"net/http"
	"net/url"
	"os"

	_ "iptv/sources/iptvorg"
	_ "iptv/sources/static"
)

func main() {
	config.LoadConfig("config.json")
	if config.Config.Proxy != "" {
		proxyUrl := must.Must2(url.Parse(config.Config.Proxy)).(*url.URL)
		http.DefaultClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	must.Must(os.MkdirAll("result", 0777))
	for _, source := range config.Config.Sources {
		must.Must(os.MkdirAll(fmt.Sprintf("result/sources/%s", source.Name), 0777))
		if source.Dedup {
			must.Must(os.MkdirAll(fmt.Sprintf("result/sources/%s/dedup", source.Name), 0777))
		}
		cha := sources.GetSource(source.Name)
		links, tinfos, err := cha.Download(source)
		must.Must(err)

		grs, infos, err := cha.Field(source, links, tinfos)
		must.Must(epg.EPGSAVE(must.Must2(epg.EPGTOEPGXML(infos)).(epg.EPGXML), fmt.Sprintf("result/sources/%s/%s.xml", source.Name, "all")))
		must.Must(err)
		for _, gr := range grs {
			must.Must(playlist.M3UGroupSaveFile([]playlist.PlayList{gr}, fmt.Sprintf("result/sources/%s/%s.m3u", source.Name, gr.Name), false))
		}
		must.Must(playlist.M3UGroupSaveFile(grs, fmt.Sprintf("result/sources/%s/%s.m3u", source.Name, "all"), false))
		if source.Dedup {
			must.Must(
				epg.EPGSAVE(
					must.Must2(
						epg.EPGTOEPGXML(infos)).(epg.EPGXML),
					fmt.Sprintf("result/sources/%s/%s/%s.xml", source.Name, "dedup", "all")))
			for _, gr := range grs {
				must.Must(playlist.M3UGroupSaveFile([]playlist.PlayList{gr}, fmt.Sprintf("result/sources/%s/%s/%s.m3u", source.Name, "dedup", gr.Name), true))
			}
			must.Must(playlist.M3UGroupSaveFile(grs, fmt.Sprintf("result/sources/%s/%s/%s.m3u", source.Name, "dedup", "all"), true))
		}
	}
}
