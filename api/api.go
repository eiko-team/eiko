package api //import "github.com/eiko-team/eiko/api"

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/eiko-team/eiko/api/consumables"
	"github.com/eiko-team/eiko/api/global"
	"github.com/eiko-team/eiko/api/list"
	"github.com/eiko-team/eiko/api/store"
	"github.com/eiko-team/eiko/api/umanagement"
	"github.com/eiko-team/eiko/api/verify"
	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/files"
	"github.com/eiko-team/eiko/misc/log"
	"github.com/eiko-team/eiko/misc/misc"
	"github.com/eiko-team/eiko/misc/structures"

	"github.com/julienschmidt/httprouter"
)

// Func is used to call the Function with the wrapper
type Func struct {
	// Function is the function wrapped
	Function func(data.Data, *http.Request) (string, error)
	// Path is the path on with you want to call the function from the api
	Path string
}

// NewFunc return a new Func struct
func NewFunc(function func(data.Data, *http.Request) (string, error), path string) Func {
	return Func{function, path}
}

// File is used to link the special file path with the URL to serve
type File struct {
	// Path of the file to serve
	Path string
	// CType AKA Content-type
	CType string
	// URLs of the served file
	URL []string
	// Title page title
	Title string
	//  User Current user
	User structures.User
}

// NewFile return a new File struct
func NewFile(path, ct, title string, url []string) File {
	return File{path, ct, url, title, structures.User{}}
}

// NewSFile return a new spacial File struct
func NewSFile(path, ct string, url []string) File {
	return File{path, ct, url, "", structures.User{}}
}

var (
	// Path of the api to serve files
	Path = os.Getenv("STATIC_PWD")

	// D api to communicate with datasctore
	D data.Data

	// Logger used to log output
	Logger = log.New(os.Stdout, "Api: ", log.Ldate|log.Ltime|log.Lshortfile)

	// FunctionsWithToken List all api functions that does require a token
	FunctionsWithToken = []Func{
		NewFunc(umanagement.UpdateToken, "/updatetoken"),
		NewFunc(umanagement.Delete, "/delete"),
		NewFunc(store.AddStore, "/store/add"),
		NewFunc(store.GetStore, "/store/get"),
		NewFunc(store.ScoreStore, "/store/score"),
		NewFunc(consumables.Store, "/consumable/add"),
		NewFunc(consumables.Get, "/consumable/get"),
		NewFunc(list.AddList, "/list/create"),
		NewFunc(list.GetLists, "/list/getall"),
		NewFunc(list.GetListContent, "/list/getcontent"),
		NewFunc(list.AddPersonnal, "/list/add/personnal"),
	}
	// Functions List all api functions that does not require a token
	Functions = []Func{
		NewFunc(global.Log, "/log"),
		NewFunc(umanagement.Login, "/login"),
		NewFunc(umanagement.Register, "/register"),
		NewFunc(verify.Email, "/verify/email"),
		NewFunc(verify.Password, "/verify/password"),
	}

	// SFiles is stored informations on special files
	SFiles = []File{
		NewSFile("js/eiko/eiko-sw.js", "application/javascript", []string{"/eiko-sw.js"}),
		NewSFile("img/EIKO.ico", "image/vnd.microsoft.icon", []string{"/favicon.ico", "/EIKO.ico"}),
		NewSFile("json/manifest.json", "application/json", []string{"/manifest.json", "/json/manifest.json"}),
		NewSFile("json/autocomplete_data.json", "application/json", []string{"/json/autocomplete_data.json"}),
	}

	// TFiles is an array of Templated files
	TFiles = []File{
		NewFile("html/eiko.html", "text/html", "Acceuil", []string{"/eiko.html", "/", "/index.html", "/l/:id"}),
		NewFile("html/login.html", "text/html", "Connexion", []string{"/login.html"}),
		NewFile("html/search.html", "text/html", "Recherche", []string{"/search/", "/search.html"}),
	}
)

// SpecialFiles simple html file server
func (file File) SpecialFiles(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	misc.LogRequest(r)

	fileContent, err := files.GetFileContent(Path + "/static/" + file.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		fmt.Fprintln(w, "{\"error\":\"invalid_file\"}")
		return
	}

	// WriteHeader MUST be called after Headers are set
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", file.CType))
	w.Header().Set("Cache-Control", "public, max-age=7776000")

	w.WriteHeader(200)
	fmt.Fprintln(w, fileContent)
}

// TemplatedFilesWrapped simple html file server
func (file File) TemplatedFilesWrapped(w http.ResponseWriter,
	r *http.Request) error {
	misc.LogRequest(r)
	templates, err := template.New("webpage").
		ParseGlob(Path + "/static/html/partial/*.html")
	if err != nil {
		return fmt.Errorf("{\"error\":\"partial_parse\"}")
	}

	fileContent, err := files.GetFileContent(Path + "/static/" + file.Path)
	if err != nil {
		return fmt.Errorf("{\"error\":\"invalid_file\"}")
	}

	h, err := templates.Parse(fileContent)
	if err != nil {
		return fmt.Errorf("{\"error\":\"parse_failed\"}")
	}

	file.User = D.User

	w.Header().Set("Content-Type", file.CType)
	err = h.Execute(w, file)
	if err != nil {
		return fmt.Errorf("{\"error\":\"could_not_exec\"}")
	}
	return nil
}

// TemplatedFiles simple html file server (wapper for TemplatedFilesWrapped)
func (file File) TemplatedFiles(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {

	if err := file.TemplatedFilesWrapped(w, r); err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w)
	}
}

// ServeFilesCustom TODO
func ServeFilesCustom(router *httprouter.Router, urlpath, filepath string) {
	fileServer := http.FileServer(http.Dir(filepath))
	router.GET(urlpath+"/*filepath",
		func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			misc.LogRequest(r)
			w.Header().Set("Vary", "Accept-Encoding")
			w.Header().Set("Cache-Control", "public, max-age=7776000")
			r.URL.Path = ps.ByName("filepath")
			fileServer.ServeHTTP(w, r)
		})
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
	for _, file := range TFiles {
		for _, URL := range file.URL {
			r.GET(URL, file.TemplatedFiles)
		}
	}

	fileType := Path + "/static/"
	fileTypeStr := "FILE_TYPE"
	if os.Getenv(fileTypeStr) == "" {
		Logger.Println(fileTypeStr + " not set using default (none)")
	} else {
		fileType += "min/"
	}
	for _, tt := range []string{"img", "js", "css"} {
		ServeFilesCustom(r, "/"+tt, fileType+tt)
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
		return
	}
	w.WriteHeader(200)
	fmt.Fprintln(w, data)
}

// WrapperFunctionCookie allows us to call the functions with rights args.
// Db must be set already.
// Read Token cookie and set the User value of D
func (fun Func) WrapperFunctionCookie(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	token, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "{\"error\":\"no_token_found\"}")
		return
	}

	D.User, err = misc.TokenToUser(D, token.Value)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "{\"error\":\"token_invalid\"}")
		return
	}
	fun.WrapperFunction(w, r, nil)
}

// InitAPI Execute the api and return the bdd configured.
func InitAPI() *httprouter.Router {
	r := httprouter.New()
	ServeFiles(r)
	for _, tt := range Functions {
		r.POST(fmt.Sprintf("/api%s", tt.Path), tt.WrapperFunction)
	}
	for _, tt := range FunctionsWithToken {
		r.POST(fmt.Sprintf("/api%s", tt.Path), tt.WrapperFunctionCookie)
	}
	return r
}

// ExecuteAPI Execute the api and return the bdd configured.
func ExecuteAPI() *httprouter.Router {

	projIDStr := "PROJECT_ID"
	projID := os.Getenv(projIDStr)
	if projID == "" {
		Logger.Fatalf("please set: '%s'", projIDStr)
	}

	D = data.InitData(projID)
	return InitAPI()
}
