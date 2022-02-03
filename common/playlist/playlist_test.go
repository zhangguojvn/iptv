package playlist_test

import (
	"iptv/common/playlist"
	"reflect"
	"testing"
)

var result = []playlist.Channel{
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "1"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "2"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "3"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "4"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "5"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "6"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "7"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "8"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "9"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "10"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "11"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "12"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "13"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "14"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "15"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "16"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "a"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "b"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "c"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "d"}},
}

var test = []playlist.Channel{
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "9"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "7"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "8"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "6"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "5"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "4"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "3"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "2"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "1"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "16"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "15"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "13"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "14"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "12"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "11"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "10"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "a"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "d"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "c"}},
	{Tvg: struct {
		Id   string "json:\"id\""
		Name string "json:\"name\""
		Url  string "json:\"url\""
	}{Id: "b"}},
}

func TestChannelSort(t *testing.T) {
	runResult := playlist.ChannelSort(test)
	if !reflect.DeepEqual(runResult, result) {
		panic("test error")
	}
}
