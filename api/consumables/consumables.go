package consumables

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/log"
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

	i.Created = time.Now()
	i.Creator = d.User.ID
	i.Source = "API"

	keyID, err := d.StoreConsumable(misc.NormalizeConsumable(i))
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

	consu, err := d.GetConsumablesTmp(misc.NormalizeQuery(i))
	if err != nil {
		Logger.Println(err)
		return "", errors.New("3.1.1")
	}

	j, err := json.Marshal(consu)
	return string(j), err
}

func GetByID(d data.Data, r *http.Request) (string, error) {
	var i structures.List
	err := misc.ParseJSON(r, &i)
	if err != nil {
		Logger.Println(err)
		return "", errors.New("3.2.0")
	}

	c, err := d.GetOneConsumable(i.ID)
	if err != nil {
		Logger.Println(err)
		return "", errors.New("3.2.1")
	}

	j, err := json.Marshal(c)
	return string(j), err
}