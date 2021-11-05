package iptv_test

import (
	"iptv/iptv"
	"net/http"
	"net/url"
	"testing"
)

func TestIPTVField(t *testing.T) {
	proxyUrl, err := url.Parse("http://127.0.0.1:9000")
	if err != nil {
		t.Error(err)
	}
	http.DefaultClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	iptv.IPTVField(iptv.IPTV{Name: "CCTV-", UrlType: "domain"}, "https://iptv-org.github.io/iptv/channels.json", false, 10)
}

func TestTestIPTV(t *testing.T) {
	proxyUrl, err := url.Parse("http://127.0.0.1:9000")
	if err != nil {
		t.Error(err)
	}
	http.DefaultClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	iptv.TestIPTV("http://116.199.5.52:8114/00000000/index.m3u8?Fsv_CMSID=&Fsv_SV_PARAM1=0&Fsv_ShiftEnable=0&Fsv_ShiftTsp=0&Fsv_chan_hls_se_idx=54&Fsv_cid=0&Fsv_ctype=LIVES&Fsv_ctype=LIVES&Fsv_filetype=1&Fsv_otype=1&Fsv_otype=1&Fsv_rate_id=0&FvSeid=5abd1660af1babb4&Pcontent_id=&Provider_id=", true, 10)
	iptv.TestIPTV("http://127.0.0.1", true, 10)
}
