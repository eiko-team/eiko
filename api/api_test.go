package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	ServeFiles(router)

	tests := []struct {
		name        string
		method      string
		URL         string
		code        int
		invalid     bool
		contentType string
	}{
		{"root", "GET", "/", http.StatusOK, false, "text/html"},
		{"index", "GET", "/index.html", http.StatusOK, false, "text/html"},
		{"login", "GET", "/login.html", http.StatusOK, false, "text/html; charset=utf-8"},
		{"app", "GET", "/eiko.html", http.StatusOK, false, "text/html"},
		{"service_worker", "GET", "/eiko-sw.js", http.StatusOK, false, "application/javascript; charset=utf-8"},
		{"favicon_Eiko", "GET", "/EIKO.ico", http.StatusOK, false, "image/vnd.microsoft.icon; charset=utf-8"},
		{"favicon", "GET", "/favicon.ico", http.StatusOK, false, "image/vnd.microsoft.icon; charset=utf-8"},
		{"manifest", "GET", "/manifest.json", http.StatusOK, false, "application/json; charset=utf-8"},
		{"js lib", "GET", "/js/lib.js", http.StatusOK, false, "application/javascript"},
		{"js color", "GET", "/js/color.js", http.StatusOK, false, "application/javascript"},
		{"not existing", "GET", "/blabla", 404, false, "text/plain; charset=utf-8"},
		{"no path", "GET", "/manifest.json", 500, true, "application/json"},
		{"no path", "GET", "/login.html", 500, true, "application/json"},
		{"no path", "GET", "/index.html", 500, true, "application/json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "no path" {
				Path = ""
			}
			t.Logf("%s", Path)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.URL, nil)
			router.ServeHTTP(w, req)
			if w.Code != tt.code {
				t.Errorf("Server error = %d, want %d for %s",
					w.Code, tt.code, tt.URL)
			}
			body := w.Body.String()
			if w.Header()["Content-Type"][0] != tt.contentType {
				t.Errorf("Content-Type = '%s' want '%s'",
					w.Header()["Content-Type"][0], tt.contentType)
			}
			if (body == "{\"error\":\"invalid_file\"}\n") != tt.invalid {
				t.Errorf("%+v", body)
			}
		})
	}
}

func TestWrapperFunction(t *testing.T) {
	router := InitAPI()
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
	router := InitAPI()
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
	router := InitAPI()

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
		{"wrong id", token, 500, "/l/test", "{\"error\":\"no_page_found\"}\n"},
		{"wrong url", token, 404, "/l", "404 page not found\n"},
		{"wrong data", token, 500, "/l/424242", "{\"error\":\"no_list_found\"}\n"},
		{"wrong path", token, 500, "/l/424242", "{\"error\":\"invalid_file\"}\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Path = os.Getenv("STATIC_PWD")
			if tt.name == "wrong data" {
				data.Error = data.ErrTest
			}
			if tt.name == "wrong path" {
				Path = ""
			}
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
			data.Error = nil
		})
	}
}

func TestExecuteAPI(t *testing.T) {
	os.Setenv("PROJECT_ID", "api_test.go")
	ExecuteAPI()
	os.Setenv("FILE_TYPE", "test")
	ExecuteAPI()
}
