package ifengcom_test

import (
	"iptv/common/config"
	"iptv/sources/ifengcom"
	"testing"
)

func TestCrawler(t *testing.T) {
	ifeng := ifengcom.IfengCom{}
	ifeng.Crawler(config.ChannelSourceConfig{Url: "https://phtv.ifeng.com/live?channel=phtvChinese"})
}

func TestGetM3uFromJS(t *testing.T) {
	ifeng := ifengcom.IfengCom{}
	ifeng.GetM3uFromJS(`q={phtvHK:"http://zb.ios.ifeng.com/live/05QGDOFQJ28/index.m3u8",phtvNews:"http://playtv-live.ifeng.com/live/06OLEEWQKN4.m3u8",phtvChinese:"http://playtv-live.ifeng.com/live/06OLEGEGM4G.m3u8"}`, "^.*:\"http://.*\\.m3u8\"$")

}
