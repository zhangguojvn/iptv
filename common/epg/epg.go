package epg

import (
	"encoding/xml"
	"io/ioutil"
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
