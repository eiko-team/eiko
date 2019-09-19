package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"eiko/api"

	"github.com/julienschmidt/httprouter"
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
	api.ServeFiles(router)

	tests := []struct {
		name   string
		method string
		URL    string
		code   int
	}{
		{"root", "GET", "/", http.StatusOK},
		{"index", "GET", "/index.html", http.StatusOK},
		{"index", "GET", "/login.html", http.StatusOK},
		{"app", "GET", "/eiko.html", http.StatusOK},
		{"service_worker", "GET", "/eiko-sw.js", http.StatusOK},
		{"favicon_Eiko", "GET", "/EIKO.ico", http.StatusOK},
		{"favicon", "GET", "/favicon.ico", http.StatusOK},
		{"manifest", "GET", "/manifest.json", http.StatusOK},
		{"not existing", "GET", "/blabla", 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.URL, nil)
			router.ServeHTTP(w, req)
			if w.Code != tt.code {
				t.Errorf("Server error = %d, want %d for %s",
					w.Code, tt.code, tt.URL)
			}
		})
	}
}
