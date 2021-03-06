package store_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/eiko-team/eiko/api/store"
	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/misc"
)

var (
	d data.Data

	token, _ = misc.UserToToken(data.UserTest)
)

func TestAddStore(t *testing.T) {
	j, _ := json.Marshal(data.StoreTest)
	s := string(j)
	tests := []struct {
		name    string
		want    string
		wantErr bool
		body    string
		useData bool
	}{
		{"sanity", `{"done":"true"}`, false, "{\"store\":" + s + "}", true},
		{"wrong json", `4.0.0`, true, "test", false},
		{"wrong data", `4.0.1`, true, "{\"store\":" + s + "}", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "wrong data" {
				data.Error = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/store/add",
				strings.NewReader(tt.body))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := store.AddStore(d, req)
			t.Logf("%v", got)
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
			if data.StoreStore != tt.useData {
				t.Errorf("data.StoreStore was no used")
			}
			data.StoreStore = false
			data.Error = nil
		})
	}
}

func TestGetStore(t *testing.T) {
	j, _ := json.Marshal(data.StoreTest)
	s := string(j)
	tests := []struct {
		name    string
		want    string
		wantErr bool
		body    string
		useData bool
	}{
		{"sanity", data.StoreRe, false, "{\"store\":" + s + "}", true},
		{"wrong json", "4.1.0", true, "}", false},
		{"wrong data", "4.1.1", true, "{\"store\":" + s + "}", true},
	}
	for _, tt := range tests {
		if tt.name == "wrong data" {
			data.Error = data.ErrTest
		}
		t.Run(tt.name, func(t *testing.T) {
			data.Store = data.StoreTest
			req, _ := http.NewRequest("POST", "/store/get",
				strings.NewReader(tt.body))
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
			if data.GetStore != tt.useData {
				t.Errorf("data.GetStore was no used")
			}
			data.GetStore = false
			data.Error = nil
		})
	}
}

func TestScoreStore(t *testing.T) {
	j, _ := json.Marshal(data.StoreTest)
	s := string(j)
	tests := []struct {
		name    string
		want    string
		wantErr bool
		body    string
		useData bool
	}{
		{"sanity", `{"done":"true"}`, false, "{\"store\":" + s + "}", true},
		{"wrong json", "4.2.0", true, "}", false},
		{"wrong data", data.ErrTest.Error(), true, "{\"store\":" + s + "}", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "wrong data" {
				data.Error2 = data.ErrTest
			}
			data.Store = data.StoreTest
			req, _ := http.NewRequest("POST", "/store/score",
				strings.NewReader(tt.body))
			got, err := store.ScoreStore(d, req)
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
			if data.ScoreStore != tt.useData {
				t.Errorf("data.ScoreStore was no used")
			}
			data.ScoreStore = false
			data.Error = nil
		})
	}
}

func TestDeleteStore(t *testing.T) {
	j, _ := json.Marshal(data.StoreTest)
	s := string(j)
	tests := []struct {
		name    string
		want    string
		wantErr bool
		body    string
		del     bool
	}{
		{"sanity", `{"done":"true"}`, false, "{\"store\":" + s + "}", true},
		{"wrong data", data.ErrTest.Error(), true, "{\"store\":" + s + "}", true},
		{"no body", `4.3.0`, true, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "wrong data" {
				data.Error = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/store/delete",
				strings.NewReader(tt.body))
			got, err := store.DeleteStore(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("DeleteStore() = '%v', want %v", got, tt.want)
			}
			if data.DeleteStore != tt.del {
				t.Errorf("data.DeleteStore = %t, want %t",
					data.DeleteStore, tt.del)
			}
			data.DeleteStore = false
			data.Error = nil
		})
	}
}

func TestUpdateStore(t *testing.T) {
	j, _ := json.Marshal(data.StoreTest)
	s := string(j)
	tests := []struct {
		name    string
		want    string
		wantErr bool
		body    string
		del     bool
	}{
		{"sanity", `{"done":"true"}`, false, s, true},
		{"wrong data", data.ErrTest.Error(), true, s, true},
		{"no body", `4.4.0`, true, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "wrong data" {
				data.Error = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/store/update",
				strings.NewReader(tt.body))
			got, err := store.UpdateStore(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("UpdateStore() = '%v', want %v", got, tt.want)
			}
			if data.UpdateStore != tt.del {
				t.Errorf("data.UpdateStore = %t, want %t",
					data.UpdateStore, tt.del)
			}
			data.UpdateStore = false
			data.Error = nil
		})
	}
}
