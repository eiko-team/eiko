package misc_test

import (
	"eiko/misc/data"
	"eiko/misc/misc"
	"eiko/misc/structures"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseJSON(t *testing.T) {
	t.Run("no request", func(t *testing.T) {
		var consumable structures.Consumable
		misc.ParseJSON(nil, &consumable)
	})
	t.Run("no body", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/test", nil)
		var consumable structures.Consumable
		misc.ParseJSON(r, &consumable)
	})
	t.Run("ParseJSON", func(t *testing.T) {
		body, _ := json.Marshal(data.ConsumableTest)
		r, _ := http.NewRequest("POST", "/test",
			strings.NewReader(string(body)))
		var consumable structures.Consumable
		if err := misc.ParseJSON(r, &consumable); err != nil {
			t.Errorf("ParseJSON() = %+v", err)
		}
		/*		if reflect.DeepEqual(data.ConsumableTest, consumable) {
					t.Errorf("ParseJSON(%+v) != %+v", consumable, data.ConsumableTest)
				}
		*/
		// TODO assert deepequals data.ConsumableTest and consumable
	})
	t.Run("Bad Json", func(t *testing.T) {
		body := "xoto"
		r, _ := http.NewRequest("POST", "/test",
			strings.NewReader(body))
		var consumable structures.Consumable
		err := misc.ParseJSON(r, &consumable)
		if err == nil {
			t.Errorf("ParseJSON(%+v) != %v", consumable, true)
		}
		if reflect.DeepEqual(data.ConsumableTest, consumable) {
			t.Errorf("ParseJSON(%+v) != %+v", consumable, data.ConsumableTest)
		}

	})
}

func TestLogRequest(t *testing.T) {
	t.Run("LogRequest", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/test", nil)
		misc.LogRequest(r)
	})
	t.Run("nil request", func(t *testing.T) {
		misc.LogRequest(nil)
	})
}

func TestToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
		email string
		pass  string
		want  bool
		err   bool
	}{
		{"sanity", "", "email", "pass", true, false},
		{"sanity", "", "email@email.em", "password", true, false},
		{"fake token", "fake.token.test", "email", "pass", false, true},
		{"invalid token", "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.VFb0qJ1LRg_4ujbZoRMXnVkUgiuKq5KxWqNdbKq_G9Vvz-S1zZa9LPxtHWKa64zDl2ofkT8F6jBt_K4riU-fPg", "email", "pass", false, true},
		{"wrong signing method", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.POstGetfAytaZS82wHcjoTyoqhMyxXiWdR7Nn7A29DNSl0EiXLdwJ6xC6AfgZWF1bOsS_TuYI3OG85AmiExREkrS6tDfTQ2B3WXlrr-wp5AokiRbz3_oB4OxG-W9KcEEbDRcZc0nH3L7LzYptiy1PtAylQGxHTWZXtGz4ht0bAecBgmpdgXMguEIcoqPJ1n3pIWk_dUZegpqx0Lka21H6XxUTxiy8OcaarA8zdnPUnV6AmNP3ecFawIFYdvJB_cm-GvpCSbr8G8y_Mllj8f4x9nBH8pQux89_6gUY618iYv7tuPWBFfEbLxtF2pZS6YC1aSfLQxeNe8djT9YjpvRZA", "email", "pass", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user structures.User
			if tt.token == "" {
				user = structures.User{
					Email:     tt.email,
					Pass:      tt.pass,
					Created:   time.Now(),
					Validated: false,
				}
				token, err := misc.UserToToken(user)
				if (err != nil) != tt.err {
					t.Errorf("UserToToken()'s error = %v", err)
				}
				tt.token = token
			}
			if got := misc.ValidateToken(tt.token); got != tt.want {
				t.Errorf("ValidateToken() = %v, want %v", got, tt)
			}
			got, err := misc.TokenToUser(tt.token)
			if (err != nil) == tt.want {
				t.Errorf("TokenToUser()'s error = %v want %v", err, tt.want)
			}
			if user.Email != got.Email {
				t.Errorf("TokenToUser() = %v, want %v", got, user)
			}
		})
	}
}

var (
	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

func TestAtoi(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{"zero", "0", 0, false},
		{"neg one", "-1", -1, false},
		{"max int", "9223372036854775807", MaxInt, false},
		{"min int", "-9223372036854775808", MinInt, false},
		{"NaN #1", "test", 0, true},
		{"NaN #2", "0x64", 0, true},
		{"datastore id#1", "5962535197259776", 5962535197259776, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := misc.Atoi(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Atoi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Atoi() = %v, want %v", got, tt.want)
			}
		})
	}
}

type struct1 struct {
	Name    string `json:"name" firestore:"Name"`
	Address string `json:"address" firestore:"Address"`
	Country string `json:"country" firestore:"Country"`
}

type struct2WithInt struct {
	Name string `json:"name" firestore:"Name"`
	Code int    `json:"code" firestore:"Code"`
}

var (
	s1 = struct1{"test name", "test address", "test country"}
	s2 = struct2WithInt{"test name", 10}
)
