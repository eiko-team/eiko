package global

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"eiko/misc/data"
	"eiko/misc/misc"
	"eiko/misc/structures"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "verify: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

// Log Logs data
func Log(d data.Data, r *http.Request) (string, error) {
	var i structures.Logging
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("3.0.0")
	}
	user, _ := misc.TokenToUser(i.Token)
	err = d.Log(structures.Log{
		Email:   user.Email,
		Log:     i.Log,
		Created: time.Now(),
	})
	if err == nil {
		return "", errors.New("3.0.1")
	}

	return "{\"done\":\"true\"}", nil
}
