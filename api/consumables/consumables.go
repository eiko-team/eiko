package consumables

import (
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
	Logger = log.New(os.Stdout, "consumables: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

// Store Stores consumable data
func Store(d data.Data, r *http.Request) (string, error) {
	var i structures.Consumable
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("3.0.0")
	}

	err = d.StoreConsumable(i)
	if err != nil {
		return "", errors.New("3.0.2")
	}

	return "{\"done\":\"true\"}", nil
}
