package global_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/eiko-team/eiko/api/global"
	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/misc"
)

var (
	d data.Data

	token, _ = misc.UserToToken(data.UserTest)
)

type Logging struct {
	Log   string `json:"message"`
	Token string `json:"user_token"`
}

func TestLog(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  bool
		wantData bool
		body     string
	}{
		{"sanity", false, true, "{\"message\":\"\",\"user_token\":\"" + token + "\"}"},
		{"no token", false, true, "{\"message\":\"\",\"user_token\":\"\"}"},
		{"wrong json", true, false, "toto"},
		{"err Log", true, true, "{\"message\":\"\",\"user_token\":\"\"}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "err Log" {
				data.Error = data.ErrTest
			}
			req, _ := http.NewRequest("POST", "/log",
				strings.NewReader(tt.body))
			got, err := global.Log(d, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Log() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != "{\"done\":\"true\"}") != tt.wantErr {
				t.Errorf("Log() = %v, wantErr %v", got, tt.wantErr)
			}
			t.Logf("Log() = %v %v", got, err)
			if data.Log != tt.wantData {
				t.Errorf("No Log")
			}
			data.Log = false
		})
	}
}
