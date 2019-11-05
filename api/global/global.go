package global

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/misc"
	"github.com/eiko-team/eiko/misc/structures"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "global: ",
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

	Logger.Printf("%v", i)
	remoteAddr := misc.SplitString(r.RemoteAddr, ":", 2)

	err = d.Log(structures.Log{
		Email:   user.Email,
		Log:     i.Log,
		Created: time.Now(),
		IP:      remoteAddr[0],
		Port:    remoteAddr[1],
	})
	if err != nil {
		return "", errors.New("3.0.1")
	}

	return "{\"done\":\"true\"}", nil
}
