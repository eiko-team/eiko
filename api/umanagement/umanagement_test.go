package umanagement_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"eiko/api/umanagement"
	"eiko/misc/data"
	"eiko/misc/structures"
)

var d data.Data

func TestManagment(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		user    structures.Login
		wantErr bool
	}{
		{"sanity", `{"token":"(.*)"}`,
			structures.Login{UserMail: "test@test.ts", UserPass: "pass"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf("{\"user_email\":\"%s\", \"user_password\":\"%s\"}",
				tt.user.UserMail, tt.user.UserPass)
			req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
			got, err := umanagement.Register(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !data.UserStored {
				t.Errorf("Data was no stored")
			}
			data.UserStored = false
			matchs := regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
			// updating token
			body = fmt.Sprintf("{\"token\":\"%s\"}", matchs[0][1])
			req, _ = http.NewRequest("POST", "/updatetoken", strings.NewReader(body))
			got, err = umanagement.UpdateToken(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			matchs = regexp.MustCompile(tt.want).FindAllStringSubmatch(got, -1)
			if len(matchs) == 0 {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
