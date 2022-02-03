package epg

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

var DownloadCache map[string]*EPGXML = make(map[string]*EPGXML)

func DownloadEPG(u string) (result EPGXML, err error) {
	if DownloadCache[u] != nil {
		result = *DownloadCache[u]
		return
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
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
	err = xml.Unmarshal(body, &result)
	if err == nil {
		DownloadCache[u] = &result
	}
	return
}
