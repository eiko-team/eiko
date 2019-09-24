package consumables

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
		return "", errors.New("3.0.1")
	}

	return "{\"done\":\"true\"}", nil
}

// Get Get some consumable data
func Get(d data.Data, r *http.Request) (string, error) {
	var i structures.Query
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("3.1.0")
	}

	consu, err := d.GetConsumable(i)
	if err != nil {
		return "", errors.New("3.1.1")
	}
	res := `{"query":[`
	if len(consu) > 0 {
		j, err := json.Marshal(consu[0])
		if err != nil {
			return "", errors.New("3.1.2")
		}
		res += string(j)
		for _, c := range consu[1:] {
			res += string(j) + ","
			j, err = json.Marshal(c)
			if err != nil {
				return "", errors.New("3.1.3")
			}
		}
	}
	return res + "]}", nil
}
