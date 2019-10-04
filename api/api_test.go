package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"eiko/misc/data"
	"eiko/misc/misc"

	"github.com/julienschmidt/httprouter"
)

var (
	token, _ = misc.UserToToken(data.UserTest)
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func TestGet400(t *testing.T) {
	router := httprouter.New()
	Path = "/home/tm/go/src/eiko"
	ServeFiles(router)

	tests := []struct {
		name    string
		method  string
		URL     string
		code    int
		invalid bool
	}{
		{"root", "GET", "/", http.StatusOK, false},
		{"index", "GET", "/index.html", http.StatusOK, false},
		{"index", "GET", "/login.html", http.StatusOK, false},
		{"app", "GET", "/eiko.html", http.StatusOK, false},
		{"service_worker", "GET", "/eiko-sw.js", http.StatusOK, false},
		{"favicon_Eiko", "GET", "/EIKO.ico", http.StatusOK, false},
		{"favicon", "GET", "/favicon.ico", http.StatusOK, false},
		{"manifest", "GET", "/manifest.json", http.StatusOK, false},
		{"js lib", "GET", "/js/lib.js", http.StatusOK, false},
		{"js color", "GET", "/js/color.js", http.StatusOK, false},
		{"not existing", "GET", "/blabla", 404, false},
		{"no path", "GET", "/manifest.json", http.StatusOK, true},
		{"no path", "GET", "/login.html", http.StatusOK, true},
		{"no path", "GET", "/index.html", http.StatusOK, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "no path" {
				Path = ""
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.URL, nil)
			router.ServeHTTP(w, req)
			if w.Code != tt.code {
				t.Errorf("Server error = %d, want %d for %s",
					w.Code, tt.code, tt.URL)
			}
			body := w.Body.String()
			if (body == "{\"error\":\"invalid_file\"}\n") != tt.invalid {
				t.Errorf("%+v", body)
			}
		})
	}
}

func TestWrapperFunction(t *testing.T) {
	router := InitApi()
	Path = "/home/tm/go/src/eiko"
	fun := Functions[0]
	t.Run("TestWrapperFunction without body", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/api"+fun.Path, nil)
		router.ServeHTTP(w, r)
		if w.Code != 500 {
			t.Errorf("TestWrapperFunction = %d, want %d",
				w.Code, 500)
		}
	})
	t.Run("TestWrapperFunction with body", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := "{\"message\":\"message\"}"
		r, _ := http.NewRequest("POST", "/api"+fun.Path,
			strings.NewReader(body))
		router.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("TestWrapperFunction = %d, want %d",
				w.Code, http.StatusOK)
		}
	})
}

func TestWrapperFunctionCookie(t *testing.T) {
	router := InitApi()
	Path = "/home/tm/go/src/eiko"
	fun := FunctionsWithToken[0]

	tests := []struct {
		name  string
		token string
		code  int
		err   string
	}{
		{"good token", token, 200, ""},
		{"no token", "", 500, "no_token_found"},
		{"invalid token", " ", 500, "token_invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/api"+fun.Path, nil)
			if tt.token != "" {
				r.Header.Set("Cookie", "token="+tt.token)
			}
			misc.LogRequest(r)
			router.ServeHTTP(w, r)
			if w.Code != tt.code {
				t.Errorf("WrapperFunctionCookie = %d, want %d",
					w.Code, tt.code)
			}

			body := w.Body.String()
			err := fmt.Sprintf("{\"error\":\"%s\"}\n", tt.err)
			if (body == err) != (tt.err != "") {
				t.Errorf("WrapperFunctionCookie = '%s', don't want '%s'",
					body, err)
			}

		})
	}
}

func TestWrapperFunctionCookieParam(t *testing.T) {
	router := InitApi()
	Path = "/home/tm/go/src/eiko"

	tests := []struct {
		name  string
		token string
		code  int
		url   string
		err   string
	}{
		{"good token", token, 200, "/l/424242", ""},
		{"no token", "", 500, "/l/424242", "{\"error\":\"no_token_found\"}\n"},
		{"invalid token", " ", 500, "/l/424242", "{\"error\":\"token_invalid\"}\n"},
		{"no id", token, 301, "/l/", "<a href=\"/l\">Moved Permanently</a>.\n\n"},
		{"wrong url", token, 404, "/l", "404 page not found\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == "" {
				data.Error = nil
			}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", tt.url, nil)
			if tt.token != "" {
				r.Header.Set("Cookie", "token="+tt.token)
			}
			misc.LogRequest(r)
			router.ServeHTTP(w, r)
			if w.Code != tt.code {
				t.Errorf("WrapperFunctionCookie = %d, want %d",
					w.Code, tt.code)
			}

			body := w.Body.String()
			if (body == tt.err) != (tt.err != "") {
				t.Errorf("WrapperFunctionCookie = '%s', want '%s'",
					body, tt.err)
			}

		})
	}
}
