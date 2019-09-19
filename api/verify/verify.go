package verify

import (
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

// Password Checks the strength of the password
// where 0 <= s <= 4
func Password(d data.Data, r *http.Request) (string, error) {
	var i structures.Password
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("2.1.0")
	}

	// TODO

	return fmt.Sprintf("{\"strength\":%d}", 4), nil
}

// Token Checks if the token is valid
func Token(d data.Data, r *http.Request) (string, error) {
	var i structures.Token
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("2.2.0")
	}

	return fmt.Sprintf("{\"valid\":\"%v\"}", misc.ValidateToken(i.Token)), nil
}
