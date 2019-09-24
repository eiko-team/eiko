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
	consu    = structures.Consumable{}
)

func TestStoreConsumable(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"sanity", `{"done":"true"}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, _ := json.Marshal(consu)
			body := fmt.Sprintf("{\"consumable\":%s}", j)
			req, _ := http.NewRequest("POST", "/consumable/add",
				strings.NewReader(body))
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
			if data.StoreConsumable == tt.wantErr {
				t.Errorf("Data was no stored")
			}
			data.StoreConsumable = false
		})
	}
}

func TestGetConsumable(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		query   string
		wantErr bool
	}{
		{"sanity", `{"query":\[\]}`, "query", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf("{\"query\":\"%s\",\"limit\":%d}",
				tt.query, 4)
			req, _ := http.NewRequest("POST", "/consumable/get",
				strings.NewReader(body))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := consumables.Get(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			t.Logf("%+v", got)
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Get() = '%v', want %v", got, tt.want)
			}
			if data.GetConsumable == tt.wantErr {
				t.Errorf("Data was no stored")
			}
			data.GetConsumable = false
		})
	}
}
