package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"eiko/api/user_managment"
	"eiko/misc/misc"

	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"
)

// Func is used to call the Function with the wrapper
type Func struct {
	// Function is the function wrapped
	Function func(http.ResponseWriter, *http.Request, context.Context,
		*datastore.Client) (string, error)

	// Path is the path on with you want to call the function from the api
	Path string
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
		Func{Function: user_managment.Login, Path: "/login"},
		Func{Function: user_managment.Register, Path: "/register"},
	}
)

// WrapperFunction allows us to call the functions with rights args.
// Db must be set already.
func (fun Func) WrapperFunction(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	misc.LogRequest(r)
	w.Header().Set("Content-Type", "application/json")
	data, err := fun.Function(w, r, Context, Datastore)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
	} else {
		w.WriteHeader(200)
		fmt.Fprintln(w, data)
	}
}

// ExecuteAPI Execute the api and return the bdd configured.
func ExecuteAPI(r *httprouter.Router) {
	Context = context.Background()

	projID := os.Getenv("PROJECT_ID")
	if projID == "" {
		Logger.Fatal("please set: 'DATASTORE_PROJECT_ID'")
	}

	var err error
	Datastore, err = datastore.NewClient(Context, projID)
	if err != nil {
		Logger.Fatalf("Could not create datastore client: %v", err)
	}

	for _, tt := range Functions {
		r.POST(tt.Path, tt.WrapperFunction)
	}
}
