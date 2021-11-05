package config

import (
	"encoding/json"
	"io/ioutil"
	"iptv/common/must"
	"os"
)

var Config config

type config struct {
	Proxy string `json:"proxy"`
	IPTV  struct {
		Source string            `json:"source"`
		Map    map[string]string `json:"map"`
		Test   struct {
			Timeout int  `json:"timeout"`
			Enable  bool `json:"enable"`
		} `json:"test"`
	} `json:"iptv"`
	EPG struct {
		Source string `json:"source"`
	} `json:"epg"`
	Logo struct {
		Source string `json:"source"`
	} `json:"logo"`
	Groups []struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		IPTV        struct {
			Regex     string `json:"regex"`
			UrlType   string `json:"url_type"`
			Languages []struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"languages"`
			Countries []struct {
				Code string `json:"code"`
				Name string `json:"name"`
			}
		} `json:"iptv"`
		EPG struct {
			Regex string `json:"regex"`
		} `json:"epg"`
	} `json:"groups"`
}

func LoadConfig(p string) {
	must.Must(json.Unmarshal(must.Must2(ioutil.ReadAll(must.Must2(os.Open(p)).(*os.File))).([]byte), &Config))

}
