package playlist

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

func TestM3u8(url string, timeout int) bool {
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
