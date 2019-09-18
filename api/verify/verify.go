package verify

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
	Logger = log.New(os.Stdout, "verify: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

// Email Checks if the email is available
func Email(d data.Data, r *http.Request) (string, error) {
	var i structures.Email
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("2.0.0")
	}

	_, err = d.GetUser(i.UserMail)
	if err == nil {
		return "", errors.New("2.0.1")
	}

	return "{\"available\":\"true\"}", nil
}
