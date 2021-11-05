package m3u

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"iptv/iptv"
)

type M3UGroup struct {
	Group string
	IPTV  []iptv.IPTV
}

func M3UGroupToFile(groups []M3UGroup, path string) error {
	cache := bytes.Buffer{}
	cache.WriteString("#EXTM3U\n")
	for _, group := range groups {
		for _, c := range group.IPTV {
			cache.WriteString(
				fmt.Sprintf(
					"#EXTINF:-1 tvg-id=\"%s\" tvg-name=\"%s\" tvg-logo=\"%s\" group-title=\"%s\",%s\n%s\n",
					c.Tvg.Id, c.Tvg.Name, c.Logo, group.Group, c.Name, c.Url,
				),
			)
		}
	}
	return ioutil.WriteFile(path, cache.Bytes(), 0777)
}
