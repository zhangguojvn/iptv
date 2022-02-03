package playlist

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Channel struct {
	Name      string   `json:"name"`
	Logo      string   `json:"logo"`
	Urls      []string `json:"url"`
	Category  string   `json:"category"`
	Languages []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"languages"`
	Countries []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"contries"`
	Tvg struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"tvg"`
}

type PlayList struct {
	Name     string
	Channels []Channel
}

func M3UGroupSaveFile(groups []PlayList, path string, isdedup bool) error {
	cache := bytes.Buffer{}
	cache.WriteString("#EXTM3U\n")
	for _, group := range groups {
		for _, c := range group.Channels {
			for _, url := range c.Urls {
				cache.WriteString(
					fmt.Sprintf(
						"#EXTINF:-1 CUID=\"%s\" tvg-id=\"%s\" tvg-chno=\"%s\" tvg-name=\"%s\" tvg-logo=\"%s\" group-title=\"%s\",%s\n%s\n",
						c.Tvg.Id, c.Tvg.Id, c.Tvg.Id, c.Tvg.Name, c.Logo, group.Name, c.Name, url,
					),
				)
				if isdedup {
					break
				}
			}
		}
	}
	return ioutil.WriteFile(path, cache.Bytes(), 0777)
}

func ChannelSort(chas []Channel) []Channel {
	var temp Channel
	for i := range chas {
		var iIsNumber = false
		iNumber, err := strconv.Atoi(chas[i].Tvg.Id)
		if err == nil {
			iIsNumber = true
		}
		for j := range chas {
			var jIsNumber = false
			jNumber, err := strconv.Atoi(chas[j].Tvg.Id)
			if err == nil {
				jIsNumber = true
			}
			if (i < j) && ((!iIsNumber && jIsNumber) || (iIsNumber && jIsNumber && iNumber > jNumber) || (!iIsNumber && !jIsNumber && chas[i].Tvg.Id > chas[j].Tvg.Id)) {
				temp = chas[i]
				chas[i] = chas[j]
				chas[j] = temp
				iNumber = jNumber
				iIsNumber = jIsNumber
			}
		}
	}
	return chas
}
