package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "Server: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
	Logger.Println("Got request!")
}

func main() {
	Logger.Println("Starting app")
	router := httprouter.New()
	router.GET("/", Index)
	Logger.Fatal(http.ListenAndServe(":80", router))
}
