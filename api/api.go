package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"eiko/api/umanagement"
	"eiko/api/verify"
	"eiko/misc/data"
	"eiko/misc/misc"

	"github.com/julienschmidt/httprouter"
)

// Func is used to call the Function with the wrapper
type Func struct {
	// Function is the function wrapped
	Function func(data.Data, *http.Request) (string, error)

	// Path is the path on with you want to call the function from the api
	Path string
}

// File is used to link the special file path with the URL to serve
type File struct {
	// Path of the file to serve
	Path string

	// Type MIME type for the file
	Type string

	// URLs of the served file
	URL []string
}

var (
	// D api to communicate with datasctore
	D data.Data

	// Logger used to log output
	Logger = log.New(os.Stdout, "Api: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Functions List all api functions
	Functions = []Func{
		{Function: umanagement.Login, Path: "/login"},
		{Function: umanagement.Register, Path: "/register"},
		{Function: umanagement.Delete, Path: "/delete"},
		{Function: umanagement.UpdateToken, Path: "/updatetoken"},
		{Function: verify.Email, Path: "/verify/email"},
		// {Function: verify.Password, Path: "/verify/password"},
		// {Function: verify.Token, Path: "/verify/token"},
	}

	// SFiles is stored informations on special files
	SFiles = []File{
		{"./static/html/index.html", "text/html", []string{"/", "/index.html"}},
		{"./static/js/eiko-sw.js", "application/x-javascript", []string{"/eiko-sw.js"}},
		{"./static/img/EIKO.ico", "image/vnd.microsoft.icon", []string{"/favicon.ico"}},
		{"./static/img/EIKO.ico", "image/vnd.microsoft.icon", []string{"/EIKO.ico"}},
		{"./static/json/manifest.json", "application/json", []string{"/manifest.json"}},
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

// ServeFiles adds to a Router special files URLs to be served. It also adds all
// static files to the Router
func ServeFiles(r *httprouter.Router) {
	Logger.Println("Adding static files")
	for _, file := range SFiles {
		for _, URL := range file.URL {
			r.GET(URL, file.SpecialFiles)
		}
	}
	for _, tt := range []string{"img", "js", "css"} {
		r.ServeFiles("/"+tt+"/*filepath", http.Dir("./static/"+tt))
	}
}

// WrapperFunction allows us to call the functions with rights args.
// Db must be set already.
func (fun Func) WrapperFunction(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	misc.LogRequest(r)
	w.Header().Set("Content-Type", "application/json")
	data, err := fun.Function(D, r)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
	} else {
		w.WriteHeader(200)
		fmt.Fprintln(w, data)
	}
}

// ExecuteAPI Execute the api and return the bdd configured.
func ExecuteAPI() *httprouter.Router {
	r := httprouter.New()
	ServeFiles(r)

	projIDStr := "PROJECT_ID"
	projID := os.Getenv(projIDStr)
	if projID == "" {
		Logger.Fatal(fmt.Sprintf("please set: '%s'", projIDStr))
	}

	D = data.InitData(projID)

	for _, tt := range Functions {
		r.POST(fmt.Sprintf("/api%s", tt.Path), tt.WrapperFunction)
	}
	return r
}
