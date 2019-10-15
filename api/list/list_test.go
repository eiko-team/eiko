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
	j, _ := json.Marshal(data.ListTest)
	l := string(j)
	tests := []struct {
		name    string
		want    string
		body    string
		wantErr bool
		useData bool
	}{
		{"sanity", data.ListRe, "{\"list\":" + l + "}", false, true},
		{"wrong json", "5.0.0", "}", true, false},
		{"wrong data", "5.0.1", "{\"list\":" + l + "}", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "wrong data" {
				data.Error = data.ErrTest
			}
			data.List = data.ListTest
			req, _ := http.NewRequest("POST", "/list/create",
				strings.NewReader(tt.body))
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
			if data.CreateList != tt.useData {
				t.Errorf("data.CreateList was not used")
			}
			data.CreateList = false
			data.Error = nil
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
		{"wrong data", "5.1.0", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "wrong data" {
				data.Error = data.ErrTest
			}
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
			if data.GetAllLists == false {
				t.Errorf("data.GetAllLists was not used")
			}
			data.GetAllLists = false
			data.Error = nil
		})
	}
}

func TestGetListContent(t *testing.T) {
	data.ListContent = data.ListContentTest
	j, _ := json.Marshal(data.ListTest)
	l := string(j)
	tests := []struct {
		name    string
		want    string
		wantErr bool
		body    string
		useData bool
	}{
		{"sanity", fmt.Sprintf("[%s]", data.ListContentRe), false,
			"{\"list\":" + l + "}", true},
		{"wrong json", "5.2.0", true, "}", false},
		{"wrong data", "5.2.1", true, "{\"list\":" + l + "}", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "wrong data" {
				data.Error = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/list/getcontent",
				strings.NewReader(tt.body))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := list.GetListContent(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListContent() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("GetListContent() = '%v', want %v", got, tt.want)
			}
			if data.GetListContent != tt.useData {
				t.Errorf("data.GetListContent was not used")
			}
			data.GetListContent = false
			data.Error = nil
		})
	}
}

func TestAddPersonnal(t *testing.T) {
	j, _ := json.Marshal(data.ListContentTest)
	data.ID = data.IDTest
	tests := []struct {
		name    string
		want    string
		wantErr bool
		body    string
	}{
		{"sanity", `{"done":true,"id":\d+}`, false, string(j)},
		{"wrong json", `5.3.0`, true, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/list/get",
				strings.NewReader(tt.body))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := list.AddPersonnal(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPersonnal() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("AddPersonnal() = '%v', want %v", got, tt.want)
			}
			if data.StoreContent == tt.wantErr {
				t.Errorf("StoreContent = %v, want = %v",
					data.StoreContent, tt.wantErr)
			}
			data.StoreContent = false
		})
	}
}
