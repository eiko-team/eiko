package verify_test

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"eiko/api/verify"
	"eiko/misc/data"
	"eiko/misc/misc"
)

var (
	d        data.Data
	token, _ = misc.UserToToken(data.UserTest)
)

func TestEmail(t *testing.T) {
	data.Error = data.ErrTest
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
		useData bool
	}{
		{"Good Email", "{\"user_email\":\"test@test.ts\"}",
			`{"available":"true"}`, false, true},
		{"wrong json", "test", `2.0.0`, true, false},
		{"wrong data", "{\"user_email\":\"test@test.ts\"}", `2.0.1`, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name != "wrong data" {
				data.Error = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/verify/email",
				strings.NewReader(tt.body))
			got, err := verify.Email(d, req)
			t.Log(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				got = err.Error()
			}
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Email() = %v(%v), want %v", got, matchs, tt.want)
			}
			if data.GetUser != tt.useData {
				t.Errorf("data.GetUser was no used")
			}
			data.GetUser = false
			data.Error = nil
		})
	}
}

func TestPassword(t *testing.T) {
	tests := []struct {
		name    string
		pass    string
		res     string
		want    bool
		wantErr bool
	}{
		{"None Password", "{\"password\":\"\"}", "0", false, false},
		{"Too Simple Password", "{\"password\":\"test\"}", "0", false, false},
		{"Simple Password", "{\"password\":\"tEst\"}", "1", false, false},
		{"Medium Password", "{\"password\":\"Test.\"}", "2", false, false},
		{"Hard Password", "{\"password\":\"tesT@test.\"}", "3", false, false},
		{"Hardest Password", "{\"password\":\"teSt@test.ts$~\"}", "4", false, false},
		{"wrong json", "test", "4", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/verify/password",
				strings.NewReader(tt.pass))
			got, err := verify.Password(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Password() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			patt := `{"strength":(\d+)}`
			if tt.wantErr {
				got = err.Error()
				patt = `2.1.0`
			}
			matchs := regexp.MustCompile(patt).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Password() = %v(%v), want %v", got, matchs, tt.want)
			}
			if !tt.wantErr && len(matchs[0]) != 0 && matchs[0][1] != tt.res {
				t.Errorf("Password() = %v, want %v", matchs[0][1], tt.res)
			}
		})
	}
}
