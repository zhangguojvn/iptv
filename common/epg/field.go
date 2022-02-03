package epg

import "regexp"

func EPGField(infos []EPG, regex string) (result []EPG, err error) {
	for _, info := range infos {
		if ok, err := regexp.Match(regex, []byte(info.Name)); !ok || err != nil {
			continue
		}
		result = append(result, info)
	}
	return
}
