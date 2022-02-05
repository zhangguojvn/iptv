package iptvorg_test

import (
	"iptv/common/config"
	"iptv/sources/iptvorg"
	"net/http"
	"net/url"
	"testing"
)

func TestDownload(t *testing.T) {
	proxyUrl, err := url.Parse("http://127.0.0.1:8889")
	if err != nil {
		t.Error(err)
	}
	http.DefaultClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	i := iptvorg.IPTVorg{}
	i.Download(config.ChannelSourceConfig{
		Name: "iptv-org",
		Url:  "https://iptv-org.github.io/iptv/channels.json",
		EPG: struct {
			Url    string                            "json:\"url\""
			Config map[string]map[string]interface{} "json:\"config\""
		}{
			Url: "http://epg.51zmt.top:8000/e.xml",
			Config: map[string]map[string]interface{}{
				"rename": {
					"CCTV1":  "CCTV-1综合",
					"CCTV2":  "CCTV-2财经",
					"CCTV3":  "CCTV-3综艺",
					"CCTV4":  "CCTV-4中文国际",
					"CCTV5":  "CCTV-5体育",
					"CCTV5+": "CCTV-5+体育赛事",
					"CCTV6":  "CCTV-6电影",
					"CCTV7":  "CCTV-7国防军事",
					"CCTV8":  "CCTV-8电视剧",
					"CCTV9":  "CCTV-9纪录",
					"CCTV10": "CCTV-10科教",
					"CCTV11": "CCTV-11戏曲",
					"CCTV12": "CCTV-12社会与法制",
					"CCTV13": "CCTV-13新闻",
					"CCTV14": "CCTV-14超清",
					"CCTV15": "CCTV-15音乐",
					"CCTV17": "CCTV-17农业农村",
				},
			},
		},
		Logo: struct {
			Url    string                            "json:\"url\""
			Config map[string]map[string]interface{} "json:\"map\""
		}{
			Url: "https://zhangguojvn.github.io/iptv/static/tvlogo/{{.Name}}.png",
			Config: map[string]map[string]interface{}{
				"link": {
					"CCTV-1综合": "https://zhangguojvn.github.io/iptv/static/tvlogo/CCTV-1综合.png",
				},
			},
		},
	})

}
