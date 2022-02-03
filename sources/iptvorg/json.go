package iptvorg

type IPTVorgJSON struct {
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	Url       string `json:"url"`
	UrlType   string `json:"-"`
	Category  string `json:"category"`
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
