package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"eiko/api/umanagement"
	"eiko/misc/misc"

	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"
)

// Func is used to call the Function with the wrapper
type Func struct {
	// Function is the function wrapped
	Function func(*http.Request, context.Context,
		*datastore.Client) (string, error)

	// Path is the path on with you want to call the function from the api
	Path string
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
	// Datastore is used to take advantage of the datastore api
	Datastore *datastore.Client

	// Context is the context of the Datastore
	Context context.Context

	// Logger used to log output
	Logger = log.New(os.Stdout, "Api: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Functions List all api functions
	Functions = []Func{
		{Function: umanagement.Login, Path: "/login"},
		{Function: umanagement.Register, Path: "/register"},
	}

	// SFiles
	SFiles = []File{
		{"./static/html/index.html", "text/html", []string{"/", "/index.html"}},
		{"./static/js/eiko-sw.js", "application/x-javascript", []string{"/eiko-sw.js"}},
		{"./static/img/EIKO.ico", "image/vnd.microsoft.icon", []string{"/favicon.ico"}},
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

// WrapperFunction allows us to call the functions with rights args.
// Db must be set already.
func (fun Func) WrapperFunction(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	misc.LogRequest(r)
	w.Header().Set("Content-Type", "application/json")
	data, err := fun.Function(r, Context, Datastore)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
	} else {
		w.WriteHeader(200)
		fmt.Fprintln(w, data)
	}
}

func ServeFiles(r *httprouter.Router) {
	for _, file := range SFiles {
		for _, URL := range file.URL {
			r.GET(URL, file.SpecialFiles)
		}
	}
	for _, tt := range []string{"img", "js", "css"} {
		r.ServeFiles("/"+tt+"/*filepath", http.Dir("./static/"+tt))
	}
}

// ExecuteAPI Execute the api and return the bdd configured.
func ExecuteAPI() *httprouter.Router {
	r := httprouter.New()

	Context = context.Background()

	projID_str := "PROJECT_ID"
	projID := os.Getenv(projID_str)
	if projID == "" {
		Logger.Fatal(fmt.Sprintf("please set: '%s'", projID_str))
	}

	var err error
	Datastore, err = datastore.NewClient(Context, projID)
	if err != nil {
		Logger.Fatalf("Could not create datastore client: %v", err)
	}

	for _, tt := range Functions {
		r.POST(fmt.Sprintf("/api%s", tt.Path), tt.WrapperFunction)
	}
	return r
}
