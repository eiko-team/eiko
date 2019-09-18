package misc_test

import (
	"testing"
	"time"

	"eiko/misc/misc"
	"eiko/misc/structures"
)

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
