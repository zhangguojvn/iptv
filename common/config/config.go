package config

import (
	"encoding/json"
	"io/ioutil"
	"iptv/common/must"
	"os"
)

var Config config

type ChannelSourceConfig struct {
	Name   string                            `json:"name"`
	Url    string                            `json:"url"`
	Config map[string]map[string]interface{} `json:"config"`
	Dedup  bool                              `json:"dedup"`
	Logo   struct {
		Url    string                            `json:"url"`
		Config map[string]map[string]interface{} `json:"map"`
	} `json:"logo"`
	Test struct {
		Timeout int  `json:"timeout"`
		Enable  bool `json:"enable"`
	} `json:"test"`
	EPG struct {
		Url    string                            `json:"url"`
		Config map[string]map[string]interface{} `json:"config"`
	} `json:"epg"`
	Groups []struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Field       struct {
			Sources struct {
				Regex      string `json:"regex"`
				UrlType    string `json:"url_type"`
				UrlKeyWord string `json:"url_keyword"`
				Languages  []struct {
					Code string `json:"code"`
					Name string `json:"name"`
				} `json:"languages"`
				Countries []struct {
					Code string `json:"code"`
					Name string `json:"name"`
				}
			} `json:"sources"`
			EPG struct {
				Regex string `json:"regex"`
			} `json:"epg"`
		} `json:"field"`
	} `json:"groups"`
}
type config struct {
	Proxy   string                `json:"proxy"`
	Sources []ChannelSourceConfig `json:"sources"`
}

func LoadConfig(p string) {
	must.Must(json.Unmarshal(must.Must2(ioutil.ReadAll(must.Must2(os.Open(p)).(*os.File))).([]byte), &Config))

}
