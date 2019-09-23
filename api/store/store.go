package store

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"eiko/misc/data"
	"eiko/misc/misc"
	"eiko/misc/structures"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "store: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

// AddStore creates a store
func AddStore(d data.Data, r *http.Request) (string, error) {
	var i structures.Store
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("4.0.0")
	}

	err = d.StoreStore(i)
	if err != nil {
		return "", errors.New("4.0.2")
	}

	return "{\"done\":\"true\"}", nil
}

// GetStore get a store
func GetStore(d data.Data, r *http.Request) (string, error) {
	var i structures.Store
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("4.1.0")
	}

	store, err := d.GetStore(i)
	if err != nil {
		return "", errors.New("4.1.2")
	}

	res, err := json.Marshal(store)
	if err != nil {
		return "", errors.New("4.1.3")
	}
	Logger.Printf("%v", string(res))
	return string(res), nil
}
