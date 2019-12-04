package store

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/log"
	"github.com/eiko-team/eiko/misc/misc"
	"github.com/eiko-team/eiko/misc/structures"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "store: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

// AddStore creates a store
// i.ID is ignored
func AddStore(d data.Data, r *http.Request) (string, error) {
	var i structures.Store
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("4.0.0")
	}

	// Ignore user input about the score
	i.Score = 0
	i.ScoreNb = 0

	err = d.StoreStore(i)
	if err != nil {
		return "", errors.New("4.0.1")
	}

	return `{"done":"true"}`, nil
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
		return "", errors.New("4.1.1")
	}

	res, err := json.Marshal(store)
	if err != nil {
		return "", errors.New("4.1.2")
	}
	return string(res), nil
}

// ScoreStore add a score to the store
func ScoreStore(d data.Data, r *http.Request) (string, error) {
	var i structures.Store
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("4.2.0")
	}
	return `{"done":"true"}`, d.ScoreStore(i)
}
