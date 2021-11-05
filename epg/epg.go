package epg

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	ID      string   `xml:"id,attr"`
	Name    struct {
		Lang string `xml:"lang,attr"`
		Text string `xml:",chardata"`
	} `xml:"display-name"`
}

type Programme struct {
	Start   string `xml:"start,attr"`
	Stop    string `xml:"stop,attr"`
	Channel string `xml:"channel,attr"`
	Title   struct {
		Lang string `xml:"lang,attr"`
		Text string `xml:",chardata"`
	} `xml:"title"`
	Desc struct {
		Lang string `xml:"lang,attr"`
		Text string `xml:",chardata"`
	} `xml:"desc"`
}

type EPGXML struct {
	XMLName   xml.Name    `xml:"tv"`
	Channel   []Channel   `xml:"channel"`
	Programme []Programme `xml:"programme"`
}

type EPG struct {
	Name      string
	Channel   Channel
	Programme []Programme
}

var EPGXMLCache map[string]EPGXML = map[string]EPGXML{}

func GetEPGXMLFromInternet(u string) (result EPGXML, err error) {
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
	return
}

func GetEPGXML(u string) (result EPGXML, err error) {
	if _, ok := EPGXMLCache[u]; !ok {
		EPGXMLCache[u], err = GetEPGXMLFromInternet(u)
		if err != nil {
			return
		}
	}
	return EPGXMLCache[u], nil
}

func EPGXMLTOEPG(e EPGXML) (result []EPG, err error) {
	var EPGMap map[string][]Programme = map[string][]Programme{}
	for _, p := range e.Programme {
		EPGMap[p.Channel] = append(EPGMap[p.Channel], p)
	}
	for _, c := range e.Channel {
		if p, ok := EPGMap[c.ID]; ok {
			result = append(result, EPG{Name: c.Name.Text, Channel: c, Programme: p})
		}
	}
	return
}

func EPGField(name string, u string) (result []EPG, err error) {
	epgxml, err := GetEPGXML(u)
	if err != nil {
		return nil, err
	}
	infos, err := EPGXMLTOEPG(epgxml)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		if ok, err := regexp.Match(name, []byte(info.Name)); !ok || err != nil {
			continue
		}
		result = append(result, info)
	}
	return
}

func EPGTOEPGXML(es []EPG) (result EPGXML, err error) {
	for _, c := range es {
		result.Channel = append(result.Channel, c.Channel)
		result.Programme = append(result.Programme, c.Programme...)
	}
	return
}

func EPGSAVE(e EPGXML, path string) error {
	b, err := xml.Marshal(e)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0777)
}
