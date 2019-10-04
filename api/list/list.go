package list

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"eiko/misc/data"
	"eiko/misc/misc"
	"eiko/misc/structures"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "list: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

// AddList creates a list
// i.ID is ignored
func AddList(d data.Data, r *http.Request) (string, error) {
	var i structures.List
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("5.0.0")
	}

	list, err := d.CreateList(i)
	if err != nil {
		return "", errors.New("5.0.1")
	}

	j, err := json.Marshal(list)
	if err != nil {
		return "", errors.New("5.0.2")
	}
	return string(j), nil
}

// GetLists gather all list for te user
func GetLists(d data.Data, r *http.Request) (string, error) {
	lists, err := d.GetAllLists()
	if err != nil {
		return "", errors.New("5.1.0")
	}

	json, err := json.Marshal(lists)
	if err != nil {
		return "", errors.New("5.1.1")
	}

	return string(json), err
}

// GetListContent gather list content for te user
func GetListContent(d data.Data, r *http.Request) (string, error) {
	var i structures.List
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("5.2.0")
	}

	content, err := d.GetListContent(i.ID)
	if err != nil {
		Logger.Println(err)
		return "", errors.New("5.2.1")
	}

	json, err := json.Marshal(content)
	if err != nil {
		return "", errors.New("5.2.2")
	}

	return string(json), err
}

// AddPersonnal add a personnal item to a list
func AddPersonnal(d data.Data, r *http.Request) (string, error) {
	var i structures.ListContent
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("5.3.0")
	}

	Logger.Printf("%+v", i)

	keyID, err := d.StoreContent(i)
	return fmt.Sprintf(`{"done":true,"id":%d}`, keyID), err
}
