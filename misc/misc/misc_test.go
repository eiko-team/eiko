package misc_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/misc"
	"github.com/eiko-team/eiko/misc/structures"
)

var (
	token, _ = misc.UserToToken(data.UserTest)
)

func ExampleParseJSON() {
	type Example struct {
		Name string `json:"name"`
	}

	req, _ := http.NewRequest("POST", "/test",
		strings.NewReader(`{"name":"test"}`))

	var example Example
	misc.ParseJSON(req, &example)

	fmt.Printf("%+v", example)

	// Output:
	// {Name:test}
}

func ExampleParseJSON_error() {
	type Example struct {
		Name string `json:"name"`
	}

	req, _ := http.NewRequest("POST", "/test",
		strings.NewReader(`bad Json`))

	var example Example
	err := misc.ParseJSON(req, &example)

	fmt.Printf("%+v\n", example)
	fmt.Printf("%s", err.Error())

	// Output:
	// {Name:}
	// invalid character 'b' looking for beginning of value
}

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
		c, _ := json.Marshal(consumable)
		d, _ := json.Marshal(data.ConsumableTest)
		cc, dd := string(c), string(d)
		if cc != dd {
			t.Errorf("ParseJSON(%s) != %s", cc, dd)
		}

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

func ExampleDumpRequest() {
	r, _ := http.NewRequest("GET", "/index.html", nil)

	fmt.Println(misc.DumpRequest(r))

	// Output:
	// GET /index.html HTTP/1.1
}

func TestDumpRequest(t *testing.T) {
	test := []struct {
		name string
		mode string
		url  string
		body string
	}{
		{"sanity", "GET", "/api/test", ""},
		{"GET with body", "GET", "/api/test", "this is the body"},
		{"POST with body", "POST", "/api/test", "this is the body"},
		{"DELETE with body", "DELETE", "/api/test", "this is the body"},
		{"UPDATE with body", "UPDATE", "/api/test", "this is the body"},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(tt.mode, tt.url,
				strings.NewReader(tt.body))
			got := misc.DumpRequest(r)
			want := fmt.Sprintf("%s %s HTTP/1.1\r\n\r\n%s", tt.mode, tt.url, tt.body)
			if got != want {
				t.Errorf("DumpRequest() = '%s' != '%s'", want, got)
			}
		})
	}
}

func ExampleLogRequest() {
	r, _ := http.NewRequest("GET", "/index.html", nil)

	misc.LogRequest(r) // Misc: 2000/01/01 00:00:00 misc.go:56: GET /index.html HTTP/1.1
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
	var D data.Data
	tests := []struct {
		name  string
		token string
		email string
		pass  string
		want  bool
		err   bool
	}{
		{"sanity", token, "email", "pass", true, false},
		{"sanity", token, "email@email.em", "password", true, false},
		{"fake token", "fake.token.test", "email", "pass", false, true},
		{"invalid token", "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.VFb0qJ1LRg_4ujbZoRMXnVkUgiuKq5KxWqNdbKq_G9Vvz-S1zZa9LPxtHWKa64zDl2ofkT8F6jBt_K4riU-fPg", "email", "pass", false, true},
		{"wrong signing method", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.POstGetfAytaZS82wHcjoTyoqhMyxXiWdR7Nn7A29DNSl0EiXLdwJ6xC6AfgZWF1bOsS_TuYI3OG85AmiExREkrS6tDfTQ2B3WXlrr-wp5AokiRbz3_oB4OxG-W9KcEEbDRcZc0nH3L7LzYptiy1PtAylQGxHTWZXtGz4ht0bAecBgmpdgXMguEIcoqPJ1n3pIWk_dUZegpqx0Lka21H6XxUTxiy8OcaarA8zdnPUnV6AmNP3ecFawIFYdvJB_cm-GvpCSbr8G8y_Mllj8f4x9nBH8pQux89_6gUY618iYv7tuPWBFfEbLxtF2pZS6YC1aSfLQxeNe8djT9YjpvRZA", "email", "pass", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.User = structures.User{
				Email:     tt.email,
				Pass:      tt.pass,
				Created:   time.Now(),
				Validated: false,
			}
			if got := misc.ValidateToken(D, tt.token); got != tt.want {
				t.Errorf("ValidateToken() = %v, want %v", got, tt)
			}
			got, err := misc.TokenToUser(D, tt.token)
			if (err != nil) == tt.want {
				t.Errorf("TokenToUser()'s error = %v want %v", err, tt.want)
			}
			if err == nil && tt.email != got.Email {
				t.Errorf("TokenToUser() = %v, want %v", got.Email, tt.email)
			}
			if data.GetUser != tt.want {
				t.Errorf("TokenToUser() = %v, want %v", data.GetUser, tt.want)
			}
			data.GetUser = false
		})
	}
}

var (
	MaxUint = ^uint(0)
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

func ExampleAtoi() {
	fmt.Println(misc.Atoi("42"))

	// Output:
	// 42 <nil>
}

func ExampleAtoi_error() {
	fmt.Println(misc.Atoi("test"))

	// Output:
	// 0 strconv.ParseInt: parsing "test": invalid syntax
}

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

func TestSplitString(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		sep    string
		lenRes int
	}{
		{"sanity", "tototata", " ", 10},
		{"url", "/api/test/:id", ":", 1},
		{"short lenRes", "/api/test/:id", ":", 0},
		{"long split", "a.b.c.d.e.f.g.h.i.j.k.l.m.o", ".", 13},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := misc.SplitString(tt.s, tt.sep, tt.lenRes)
			t.Logf("%v -> %q(%d)", tt, got, len(got))
			if len(got) < tt.lenRes {
				t.Errorf("SplitString(%q) = %q, want %d", tt, got, tt.lenRes)
			}
		})
	}
}

func TestNormalizeName(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		val  string
		want string
	}{
		{"sanity", "sanity", "sanity"},
		{"simple", "SanitY", "sanity"},
		{"with space", "San%20itY", "san ity"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := misc.NormalizeName(tt.val); got != tt.want {
				t.Errorf("NormalizeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeConsumable(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		c    structures.Consumable
		want structures.Consumable
	}{
		{"sanity", structures.Consumable{Name: "Test"},
			structures.Consumable{Name: "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := misc.NormalizeConsumable(tt.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalizeConsumable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeQuery(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		q    structures.Query
		want structures.Query
	}{
		{"sanity", structures.Query{Query: "Sample", Limit: 0},
			structures.Query{Query: "sample", Limit: 20}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := misc.NormalizeQuery(tt.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalizeQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleMergeUser() {
	UserWithMissingInformations := structures.User{
		Email: "test@test.test",
		ID:    42,
	}

	SameUserWithMoreInformations := structures.User{
		Pass:      "Pass",
		Created:   time.Now(),
		Validated: true,
		ID:        21,
	}
	NewUser := misc.MergeUser(UserWithMissingInformations,
		SameUserWithMoreInformations)

	fmt.Println(NewUser.Email)
	fmt.Println(NewUser.Pass)
	fmt.Println(NewUser.Validated)
	fmt.Println(NewUser.ID)

	// Output:
	// test@test.test
	// Pass
	// true
	// 42
}

func TestMergeUser(t *testing.T) {
	fullUser := structures.User{
		Email:     "test",
		Pass:      "Pass",
		Created:   time.Now(),
		Validated: true,
		ID:        42,
	}
	tests := []struct {
		name  string
		user1 structures.User
		user2 structures.User
		want  structures.User
	}{
		{"sanity", structures.User{}, structures.User{}, structures.User{}},
		{"no field", structures.User{Email: "test"}, structures.User{}, structures.User{Email: "test"}},
		{"same field", structures.User{Email: "test"}, structures.User{Email: "not test"}, structures.User{Email: "test"}},
		{"same field not set", structures.User{}, structures.User{Email: "test"}, structures.User{Email: "test"}},
		{"full fields", fullUser, structures.User{}, fullUser},
		{"none fields set", structures.User{}, fullUser, fullUser},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := misc.MergeUser(tt.user1, tt.user2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleMergeStore() {
	StoreWithMissingInformations := structures.Store{
		Name: "Store Test.",
		ID:   42,
	}

	SameStoreWithMoreInformations := structures.Store{
		Name:    "The Best Test Store",
		ID:      21,
		Country: "In",
	}
	NewStore := misc.MergeStore(StoreWithMissingInformations,
		SameStoreWithMoreInformations)

	fmt.Println(NewStore.Name)
	fmt.Println(NewStore.ID)
	fmt.Println(NewStore.Country)

	// Output:
	// The Best Test Store
	// 42
	// In
}

func TestMergeStore(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		i1   structures.Store
		i2   structures.Store
		want structures.Store
	}{
		{"sanity", structures.Store{}, structures.Store{}, structures.Store{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := misc.MergeStore(tt.i1, tt.i2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeStrings(t *testing.T) {
	tests := []struct {
		name string
		s1   []string
		s2   []string
		want []string
	}{
		{"sanity", []string{}, []string{}, []string{}},
		{"no S2", []string{"test"}, []string{}, []string{"test"}},
		{"no S1", []string{}, []string{"test"}, []string{"test"}},
		{"both filled", []string{"test s1"}, []string{"test s2"}, []string{"test s1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := misc.MergeStrings(tt.s1, tt.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeInt(t *testing.T) {
	tests := []struct {
		name string
		i1   int
		i2   int
		want int
	}{
		{"sanity", 0, 0, 0},
		{"no i1", 0, 42, 42},
		{"no i2", 42, 0, 42},
		{"both filled", 42, 21, 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := misc.MergeInt(tt.i1, tt.i2); got != tt.want {
				t.Errorf("MergeInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeFloat(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		f1   float64
		f2   float64
		want float64
	}{
		{"sanity", 0, 0, 0},
		{"no f1", 0, 0.42, 0.42},
		{"no f2", 0.42, 0, 0.42},
		{"both filled", 0.42, 0.21, 0.42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := misc.MergeFloat(tt.f1, tt.f2); got != tt.want {
				t.Errorf("MergeFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeString(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		s1   string
		s2   string
		want string
	}{
		{"sanity", "", "", ""},
		{"no S2", "test", "", "test"},
		{"no S1", "", "test", "test"},
		{"both filled", "test s1", "test s2", "test s1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := misc.MergeString(tt.s1, tt.s2); got != tt.want {
				t.Errorf("MergeString() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
