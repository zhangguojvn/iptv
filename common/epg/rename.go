package epg

func Rename(epgRAWs []EPG, config map[string]interface{}) (epgResult []EPG) {
	for _, epgRaw := range epgRAWs {
		if config != nil && config[epgRaw.Name] != nil {
			epgRaw.Channel.Name.Text = config[epgRaw.Name].(string)
			epgRaw.Name = config[epgRaw.Name].(string)
			epgResult = append(epgResult, epgRaw)
		} else {
			epgResult = append(epgResult, epgRaw)
		}
	}
	return
}
