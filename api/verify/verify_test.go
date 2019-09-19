package verify_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"eiko/api/verify"
	"eiko/misc/data"
	"eiko/misc/misc"
	"eiko/misc/structures"
)

var (
	d        data.Data
	token, _ = misc.UserToToken(structures.User{
		Email:     "test@test.ts",
		Pass:      "pass",
		Created:   time.Now(),
		Validated: false,
	})
)

func TestToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Good Token", token, false, false},
		{"Fake Token", "Fake token", false, false},
		{"Invalid Token", "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.VFb0qJ1LRg_4ujbZoRMXnVkUgiuKq5KxWqNdbKq_G9Vvz-S1zZa9LPxtHWKa64zDl2ofkT8F6jBt_K4riU-fPg", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf("{\"token\":\"%s\"}", tt.token)
			req, _ := http.NewRequest("POST", "/verify/token", strings.NewReader(body))
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
