package main

import (
	"eiko/api"
	"fmt"
	// "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"eiko/misc/misc"
	"github.com/julienschmidt/httprouter"
)

// Config server configuration struct
type Config struct {
	Port string
}

type File struct {
	// Path of the file to serve
	Path string

	// Type MIME type for the file
	Type string

	// URLs of the served file
	URL []string
}

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "Server: ", log.Ldate|log.Ltime|log.Lshortfile)
	// Configuration global configuration for the app
	Configuration = Config{}

	// SFiles
	SFiles = []File{
		File{"./static/html/index.html", "text/html", []string{"/", "/index.html"}},
		File{"./static/js/eiko-sw.js", "application/x-javascript", []string{"/eiko-sw.js"}},
		File{"./static/img/EIKO.ico", "image/vnd.microsoft.icon", []string{"/favicon.ico"}},
		File{"./static/json/manifest.json", "application/json", []string{"/manifest.json"}},
	}
)

// SpecialFiles simple html file server
func (file File) SpecialFiles(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	b, _ := ioutil.ReadFile(file.Path)
	w.Header().Set("Content-Type", file.Type)
	fmt.Fprint(w, string(b))
	misc.LogRequest(r)
}

func main() {
	// Configuration
	port := os.Getenv("PORT")
	if port == "" {
		Configuration.Port = "8080"
	} else {
		Configuration.Port = port
	}
	Logger.Printf("Port: %s\n", Configuration.Port)

	// Special Files
	router := httprouter.New()
	for _, file := range SFiles {
		for _, URL := range file.URL {
			router.GET(URL, file.SpecialFiles)
		}
	}
	for _, tt := range []string{"img", "js", "css"} {
		router.ServeFiles("/"+tt+"/*filepath", http.Dir("./static/"+tt))
	}
	Logger.Println("Starting api")
	api.ExecuteAPI(router)
	Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Configuration.Port),
		router))
}
