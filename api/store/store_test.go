package store_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"eiko/api/store"
	"eiko/misc/data"
	"eiko/misc/misc"
)

var (
	d data.Data

	token, _ = misc.UserToToken(data.UserTest)
)

func TestAddStore(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		token   string
		wantErr bool
	}{
		{"sanity", `{"done":"true"}`, token, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, _ := json.Marshal(data.StoreTest)
			body := fmt.Sprintf("{\"store\":%s}", j)
			req, _ := http.NewRequest("POST", "/store/add",
				strings.NewReader(body))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", tt.token))
			got, err := store.AddStore(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("AddStore() = '%v', want %v", got, tt.want)
			}
			if data.StoreStore == tt.wantErr {
				t.Errorf("Data was no stored")
			}
			data.StoreStore = false
		})
	}
}

func TestGetStore(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		token   string
		wantErr bool
	}{
		{"sanity", `{"name":"[a-z ]+","address":"[a-z ]+","country":"[a-z ]+","zip":"[a-z ]+","user_rating":\d+,"geohash":\d+,"ID":\d+}`, token, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.Store = data.StoreTest
			j, _ := json.Marshal(data.StoreTest)
			body := fmt.Sprintf("{\"user_token\":\"%s\",\"store\":%s}",
				tt.token, j)
			req, _ := http.NewRequest("POST", "/store/get",
				strings.NewReader(body))
			got, err := store.GetStore(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("GetStore() = '%v', want %v", got, tt.want)
			}
			if data.GetStore == tt.wantErr {
				t.Errorf("Data was no retrieved")
			}
			data.GetStore = false
		})
	}
}
