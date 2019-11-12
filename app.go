package main

import (
	"fmt"
	// "html/template"
	"net/http"
	"os"

	"github.com/eiko-team/eiko/api"
	"github.com/eiko-team/eiko/misc/log"
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
	Logger.Println("Starting api")
	Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Configuration.Port),
		api.ExecuteAPI()))
}
