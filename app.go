package main

import (
	"eiko/api"
	"fmt"
	// "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// Config server configuration struct
type Config struct {
	Port string
}

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "Server: ", log.Ldate|log.Ltime|log.Lshortfile)
	// Configuration global configuration for the app
	Configuration = Config{}
)

// Index simple html file server
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, _ := ioutil.ReadFile("./static/html/index.html")
	fmt.Fprint(w, string(b))
	Logger.Println("Accessing /")
}

// API test
func API(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
	Logger.Println("Got request for the api!")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		Configuration.Port = "8080"
	} else {
		Configuration.Port = port
	}
	Logger.Printf("Port: %s\n", Configuration.Port)
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/", API)
	for _, tt := range []string{"img", "js", "css"} {
		router.ServeFiles("/"+tt+"/*filepath", http.Dir("./static/"+tt))
	}
	Logger.Println("Starting app")
	api.ExecuteAPI(router)
	Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Configuration.Port),
		router))
}
