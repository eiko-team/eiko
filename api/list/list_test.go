package list_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"eiko/api/list"
	"eiko/misc/data"
	"eiko/misc/misc"
)

var (
	d data.Data

	token, _ = misc.UserToToken(data.UserTest)
)

func TestAddList(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"sanity", data.ListRe, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.List = data.ListTest
			j, _ := json.Marshal(data.ListTest)
			body := fmt.Sprintf("{\"list\":%s}", j)
			req, _ := http.NewRequest("POST", "/list/create",
				strings.NewReader(body))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := list.AddList(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("AddList() = '%v', want %v", got, tt.want)
			}
			if data.CreateList == tt.wantErr {
				t.Errorf("Data was no listd")
			}
			data.CreateList = false
		})
	}
}

func TestGetLists(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"sanity", fmt.Sprintf("[%s]", data.ListRe), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.List = data.ListTest
			req, _ := http.NewRequest("POST", "/list/getall", nil)
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := list.GetLists(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("GetLists() = '%v', want %v", got, tt.want)
			}
			if data.GetAllLists == tt.wantErr {
				t.Errorf("Data was no listd")
			}
			data.GetAllLists = false
		})
	}
}

func TestGetListContent(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"sanity", fmt.Sprintf("[%s]", data.ListContentRe), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.ListContent = data.ListContentTest
			j, _ := json.Marshal(data.ListTest)
			body := fmt.Sprintf("{\"list\":%s}", j)
			req, _ := http.NewRequest("POST", "/list/get",
				strings.NewReader(body))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := list.GetListContent(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("GetListContent() = '%v', want %v", got, tt.want)
			}
			if data.GetListContent == tt.wantErr {
				t.Errorf("Data was no listd")
			}
			data.GetListContent = false
		})
	}
}
