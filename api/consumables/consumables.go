package consumables

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/misc"
	"github.com/eiko-team/eiko/misc/structures"
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

	keyID, err := d.StoreConsumable(i)
	if err != nil {
		return "", errors.New("3.0.1")
	}

	return fmt.Sprintf(`{"done":true,"id":%d}`, keyID), err
}

// Get Get some consumable data
func Get(d data.Data, r *http.Request) (string, error) {
	var i structures.Query
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("3.1.0")
	}

	consu, err := d.GetConsumables(i)
	if err != nil {
		return "", errors.New("3.1.1")
	}

	j, err := json.Marshal(consu)
	return string(j), err
}
