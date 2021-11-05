package iptv

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"time"

	"log"
)

type IPTV struct {
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	Url       string `json:"url"`
	UrlType   string `json:"-"`
	Category  string `json:"category"`
	Languages []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"languages"`
	Countries []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"contries"`
	Tvg *struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"tvg"`
}

var iPTVCache map[string][]IPTV = map[string][]IPTV{}

func GetIPTVFromNetwork(u string) (result []IPTV, err error) {
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

func GetIPTV(u string) (result []IPTV, err error) {
	if iPTVCache[u] == nil {
		iPTVCache[u], err = GetIPTVFromNetwork(u)
		if err != nil {
			return nil, err
		}
	}
	return iPTVCache[u], nil
}

func IPTVField(i IPTV, u string, istest bool, timeout int) (result []IPTV, err error) {
	var resultMap map[string][]IPTV = map[string][]IPTV{}
	sources, err := GetIPTV(u)
	if err != nil {
		return
	}
	for _, source := range sources {

		if i.Name != "" {

			if ok, err := regexp.Match(i.Name, []byte(source.Name)); !ok || err != nil {
				continue
			}
		}
		if i.Languages != nil && !func() bool {
			if source.Languages == nil {
				return false
			}
			for _, requestL := range i.Languages {
				for _, haveL := range source.Languages {
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
		if i.Countries != nil && !func() bool {
			if source.Countries == nil {
				return false
			}

			for requestC := range i.Countries {
				for haveC := range source.Countries {
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
		if i.Tvg != nil && !reflect.DeepEqual(source.Tvg, i.Tvg) {
			continue
		}
		resultMap[source.Name] = append(resultMap[source.Name], source)
	}

	for tvname, iptvcs := range resultMap {
		var ishave = false
		for _, iptvc := range iptvcs {
			if !ishave {
				if i.UrlType != "" {
					if func() bool {
						switch i.UrlType {
						case "ip":
							if ok, err := regexp.Match("([0-9]{0,3}\\.){3}[0-9]{0,3}(:|/)", []byte(iptvc.Url)); !ok || err != nil {
								return false
							}
						case "domain":
							if ok, err := regexp.Match("([0-9]{0,3}\\.){3}[0-9]{0,3}(:|/)", []byte(iptvc.Url)); ok || err != nil {
								return false
							}
						default:
							return true
						}
						return true
					}() {
						if TestIPTV(iptvc.Url, istest, timeout) {
							ishave = true
							result = append(result, iptvc)
						}
					}
				} else {
					if TestIPTV(iptvc.Url, istest, timeout) {
						ishave = true
						result = append(result, iptvc)
					}
				}

			}
		}
		if !ishave && i.UrlType != "" {
			log.Println(tvname + ": Can't find prefer type url.")
			for _, iptvc := range iptvcs {
				if !ishave {
					if TestIPTV(iptvc.Url, istest, timeout) {
						result = append(result, iptvcs[0])
						ishave = true
					}
				}
			}
		}
	}
	return
}

func TestIPTV(url string, istest bool, timeout int) bool {
	if !istest {
		return true
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
	req.Header.Set("Cookie", "CONSENT=YES+cb; YSC=DwKYllHNwuw")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	if ok, err := regexp.Match("^#EXTM3U", body); ok && err == nil {
		return true
	}
	return false
}
