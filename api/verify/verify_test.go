package verify_test

import (
	"fmt"
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
	token, _ = misc.UserToToken(data.TestUser)
)

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		want    bool
		wantErr bool
	}{
		{"Good Email", "test@test.ts", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.Error = data.TestError
			body := fmt.Sprintf("{\"user_email\":\"%s\"}", tt.email)
			req, _ := http.NewRequest("POST", "/verify/email",
				strings.NewReader(body))
			got, err := verify.Email(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			matchs := regexp.MustCompile(`{"available":"true"}`).FindAllStringSubmatch(got, -1)
			if (len(matchs) == 0) != tt.want {
				t.Errorf("Email() = %v(%v), want %v", got, matchs, tt.want)
			}
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
		{"None Password", "", "0", false, false},
		{"Too Simple Password", "test", "0", false, false},
		{"Simple Password", "tEst", "1", false, false},
		{"Medium Password", "Test.", "2", false, false},
		{"Hard Password", "tesT@test.", "3", false, false},
		{"Hardest Password", "teSt@test.ts$~", "4", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf("{\"password\":\"%s\"}", tt.pass)
			req, _ := http.NewRequest("POST", "/verify/password",
				strings.NewReader(body))
			got, err := verify.Password(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Password() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			matchs := regexp.MustCompile(`{"strength":(\d+)}`).FindAllStringSubmatch(got, -1)
			if (len(matchs) == 0) != tt.want {
				t.Errorf("Password() = %v(%v), want %v", got, matchs, tt.want)
			}
			if len(matchs) != 0 && len(matchs[0]) != 0 && matchs[0][1] != tt.res {
				t.Errorf("Password() = %v, want %v", matchs[0][1], tt.res)
			}
		})
	}
}

func TestToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    bool
		wantErr bool
	}{
		{"Good Token", token, false, false},
		{"Fake Token", "Fake token", false, false},
		{"Invalid Token", "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.VFb0qJ1LRg_4ujbZoRMXnVkUgiuKq5KxWqNdbKq_G9Vvz-S1zZa9LPxtHWKa64zDl2ofkT8F6jBt_K4riU-fPg", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf("{\"token\":\"%s\"}", tt.token)
			req, _ := http.NewRequest("POST", "/verify/token",
				strings.NewReader(body))
			got, err := verify.Token(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			matchs := regexp.MustCompile(`{"valid":"(.*)"}`).FindAllStringSubmatch(got, -1)
			if (len(matchs) == 0) != tt.want {
				t.Errorf("Login() = %v(%v), want %v", got, matchs, tt.want)
			}
		})
	}
}
