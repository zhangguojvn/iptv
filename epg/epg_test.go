package epg_test

import (
	"iptv/epg"
	"net/http"
	"net/url"
	"testing"
)

func TestGetEPGXMLFromInternet(t *testing.T) {
	proxyUrl, err := url.Parse("http://127.0.0.1:9000")
	if err != nil {
		t.Error(err)
	}
	http.DefaultClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	_, err = epg.GetEPGXMLFromInternet("http://epg.51zmt.top:8000/e.xml")
	if err != nil {
		t.Error(err)
	}
}

func TestEPGXMLField(t *testing.T) {
	proxyUrl, err := url.Parse("http://127.0.0.1:9000")
	if err != nil {
		t.Error(err)
	}
	http.DefaultClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	_, err = epg.EPGField("CCTV", "http://epg.51zmt.top:8000/e.xml")
	if err != nil {
		t.Error(err)
	}
}
