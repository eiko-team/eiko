package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"eiko/api/consumables"
	"eiko/api/global"
	"eiko/api/list"
	"eiko/api/store"
	"eiko/api/umanagement"
	"eiko/api/verify"
	"eiko/misc/data"
	"eiko/misc/files"
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
	// CType AKA Content-type
	CType string
	// URLs of the served file
	URL []string
	// Title page title
	Title string
	// Type MIME type for the file
	Type string
}

// Page Page struct to contain most elements on a page
type Page struct {
	// Title page title
	Title string
	// ID List id
	ID int64
	// Name of the parameter
	Name string
	// URL of the GET router
	URL string
	// Type set usage of minimified files
	Type string
}

var (
	// D api to communicate with datasctore
	D data.Data

	// Logger used to log output
	Logger = log.New(os.Stdout, "Api: ", log.Ldate|log.Ltime|log.Lshortfile)

	// FunctionsWithToken List all api functions that does require a token
	FunctionsWithToken = []Func{
		{Function: umanagement.Delete, Path: "/delete"},
		{Function: umanagement.UpdateToken, Path: "/updatetoken"},
		{Function: store.AddStore, Path: "/store/add"},
		{Function: store.GetStore, Path: "/store/get"},
		{Function: consumables.Store, Path: "/consumable/add"},
		{Function: list.AddList, Path: "/list/create"},
		{Function: list.GetLists, Path: "/list/getall"},
		{Function: list.GetListContent, Path: "/list/get"},
	}
	// Functions List all api functions that does not require a token
	Functions = []Func{
		{Function: umanagement.Login, Path: "/login"},
		{Function: umanagement.Register, Path: "/register"},
		{Function: verify.Email, Path: "/verify/email"},
		{Function: verify.Password, Path: "/verify/password"},
		{Function: global.Log, Path: "/log"},
	}

	// SFiles is stored informations on special files
	SFiles = []File{
		{"./static/html/login.html", "text/html", []string{"/login.html"}, "", ""},
		{"./static/html/eiko.html", "text/html", []string{"/eiko.html", "/", "/index.html"}, "Acceuil", ""},
		{"./static/js/eiko/eiko-sw.js", "application/x-javascript", []string{"/eiko-sw.js"}, "", ""},
		{"./static/img/EIKO.ico", "image/vnd.microsoft.icon", []string{"/favicon.ico"}, "", ""},
		{"./static/img/EIKO.ico", "image/vnd.microsoft.icon", []string{"/EIKO.ico"}, "", ""},
		{"./static/json/manifest.json", "application/json", []string{"/manifest.json"}, "", ""},
	}

	// FunctionsWithParam slice of Page
	FunctionsWithParam = []Page{
		{"", 0, "id", "/l/:id", ""},
	}
)

// SpecialFiles simple html file server
func (file File) SpecialFiles(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	misc.LogRequest(r)
	w.Header().Set("Content-Type", "application/json")

	fileContent, err := files.GetFileContent(file.Path)
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"invalid_file\"}")
		return
	}

	h, err := template.New("webpage").Parse(fileContent)
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"parse_failed\"}")
		return
	}

	w.Header().Set("Content-Type", file.CType)
	err = h.Execute(w, file)
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"could_not_exec\"}")
		return
	}
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

// WrapperFunctionCookie allows us to call the functions with rights args.
// Db must be set already.
// Read Token cookie and set the User value of D
func (fun Func) WrapperFunctionCookie(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	token, err := r.Cookie("token")
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"no_token_found\"}")
		return
	}
	D.User, err = misc.TokenToUser(token.Value)
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"token_invalid\"}")
		return
	}
	fun.WrapperFunction(w, r, nil)
}

// WrapperFunctionCookieParam allows us to call the functions with rights args.
// Db must be set already.
// Read Token cookie and set the User value of D
func (fun Func) WrapperFunctionCookieParam(w http.ResponseWriter,
	r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	token, err := r.Cookie("token")
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"no_token_found\"}")
		return
	}
	D.User, err = misc.TokenToUser(token.Value)
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"token_invalid\"}")
		return
	}
	fun.WrapperFunction(w, r, nil)
}

// WrapperFunctionCookieParam return to the user the full page with dinamic
// content
func (page Page) WrapperFunctionCookieParam(w http.ResponseWriter,
	r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	token, err := r.Cookie("token")
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"no_token_found\"}")
		return
	}
	D.User, err = misc.TokenToUser(token.Value)
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"token_invalid\"}")
		return
	}

	pageID, err := strconv.Atoi(ps.ByName(page.Name))
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"no_page_found\"}")
		return
	}
	page.ID = int64(pageID)

	list, err := D.GetList(page.ID)
	if err != nil {
		Logger.Printf("%+v", err)
		fmt.Fprintln(w, "{\"error\":\"no_list_found\"}")
		return
	}
	page.Title = list.Name

	file, err := files.GetFileContent("./static/html/eiko.html")
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"invalid_file\"}")
		return
	}

	h, err := template.New("webpage").Parse(file)
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"parse_failed\"}")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Execute(w, page)
	if err != nil {
		fmt.Fprintln(w, "{\"error\":\"could_not_exec\"}")
		return
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
	for _, tt := range FunctionsWithToken {
		r.POST(fmt.Sprintf("/api%s", tt.Path), tt.WrapperFunctionCookie)
	}
	for _, tt := range FunctionsWithParam {
		r.GET(tt.URL, tt.WrapperFunctionCookieParam)
	}
	return r
}
