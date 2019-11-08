package consumables_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"eiko/api/consumables"
	"eiko/misc/data"
	"eiko/misc/misc"
	"eiko/misc/structures"
)

var (
	d        data.Data
	token, _ = misc.UserToToken(data.UserTest)
)

func TestStoreConsumable(t *testing.T) {
	j, _ := json.Marshal(data.ConsumableTest)
	c := string(j)
	tests := []struct {
		name    string
		want    string
		wantErr bool
		body    string
		useData bool
	}{
		{"sanity", `{"done":true,"id":\d+}`, false, "{\"consumable\":" + c + "}",
			true},
		{"bad json", `3.0.0`, true, "}", false},
		{"wrong data", `3.0.1`, true, "{\"consumable\":" + c + "}", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "wrong data" {
				data.Error = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/consumable/add",
				strings.NewReader(tt.body))
			t.Log(misc.DumpRequest(req))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := consumables.Store(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Store() = '%v', want %v", got, tt.want)
			}
			if data.StoreConsumable != tt.useData {
				t.Errorf("data.StoreConsumable not used")
			}
			data.StoreConsumable = false
			data.Error = nil
		})
	}
}

func TestGetConsumables(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		body    string
		wantErr bool
		useData bool
	}{
		{"sanity", `\[\]`, "{\"query\":\"query\"}", false, true},
		{"wrong json in body", `3.1.0`, "}", true, false},
		{"simple", fmt.Sprintf(`\[(%s,?)+\]`, data.ConsumablesRe),
			"{\"query\":\"query\"}", false, true},
		{"no data", `3.1.1`, "{\"query\":\"query\"}", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "simple" {
				data.Consumables = []structures.Consumables{
					structures.Consumables{
						Consumable: data.ConsumableTest,
						Store:      data.StoreTest,
						Stock:      structures.Stock{},
					},
					structures.Consumables{
						Consumable: data.ConsumableTest,
						Store:      data.StoreTest,
						Stock:      structures.Stock{},
					},
				}
			}
			if tt.name == "no data" {
				data.Error = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/consumable/get",
				strings.NewReader(tt.body))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := consumables.Get(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			t.Logf("test: '%+v'", tt)
			t.Logf("got: '%+v'", got)
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Get() = '%v', want something like: '%v'",
					got, tt.want)
			}
			if data.GetConsumables != tt.useData {
				t.Errorf("Data was no stored")
			}
			data.GetConsumables = false
		})
	}
}
