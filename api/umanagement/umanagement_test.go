package umanagement_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"eiko/api/umanagement"
	"eiko/misc/data"
	"eiko/misc/misc"
)

var (
	d data.Data

	token, _ = misc.UserToToken(data.UserTest)
)

func TestLogin(t *testing.T) {
	data.User = data.UserTest
	tests := []struct {
		name    string
		want    string
		body    string
		wantErr bool
		useData bool
	}{
		{"sanity", `{"token":"(.*)"}`,
			"{\"user_email\":\"test@test.ts\",\"user_password\":\"pass\"}",
			false, true},
		{"wrong json", `1.0.0`, "test", true, false},
		{"wrong data", `1.0.1`,
			"{\"user_email\":\"test@test.ts\",\"user_password\":\"pass\"}",
			true, true},
		{"wrong hash", `1.0.2`,
			"{\"user_email\":\"wrong@test.ts\",\"user_password\":\"fak pass\"}",
			true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "wrong data" {
				data.Error = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/login",
				strings.NewReader(tt.body))
			got, err := umanagement.Login(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
			if data.GetUser != tt.useData {
				t.Errorf("Data was no retrieved")
			}
			data.GetUser = false
			data.Error = nil
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name         string
		want         string
		body         string
		wantErr      bool
		useDataGet   bool
		useDataStore bool
	}{
		{"sanity", `{"token":"(.*)"}`,
			"{\"user_email\":\"test@test.ts\",\"user_password\":\"pass\"}",
			false, true, true},
		{"wrong json", `1.1.0`, "test", true, false, false},
		{"wrong data get", `1.1.1`,
			"{\"user_email\":\"test@test.ts\",\"user_password\":\"pass\"}",
			true, true, false},
		{"wrong data store", `1.1.3`,
			"{\"user_email\":\"test@test.ts\",\"user_password\":\"pass\"}",
			true, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name != "wrong data get" {
				data.Error = data.ErrTest
			}
			if tt.name == "wrong data store" {
				data.Error2 = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/register",
				strings.NewReader(tt.body))
			got, err := umanagement.Register(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Register() = %v, want %v", got, tt.want)
			}
			if data.GetUser != tt.useDataGet {
				t.Errorf("data.GetUser was no used")
			}
			data.GetUser = false
			if data.StoreUser != tt.useDataStore {
				t.Errorf("data.StoreUser was no used")
			}
			data.StoreUser = false
			data.Error = nil
			data.Error2 = nil
		})
	}
}

func TestDelete(t *testing.T) {
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
			body := fmt.Sprintf("{\"token\":\"%s\"}", tt.token)
			req, _ := http.NewRequest("POST", "/delete",
				strings.NewReader(body))
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", tt.token))

			got, err := umanagement.Delete(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Delete() = '%v', want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateToken(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"sanity", `{"token":"(.*)"}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/updatetoken", nil)
			req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
			got, err := umanagement.UpdateToken(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("UpdateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
