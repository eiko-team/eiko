package misc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"cloud.google.com/go/datastore"
	"eiko/misc/structures"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "Misc: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// UniqEmail is used to find if a email is already used in the datastore
func UniqEmail(ctx context.Context, Users, UserMail string,
	client *datastore.Client) error {
	// Finding if the email is unique
	var users []structures.User
	q := datastore.NewQuery(Users).Filter("Email =", UserMail).Limit(1)
	if _, err := client.GetAll(ctx, q, &users); err != nil {
		return err
	}
	if len(users) != 0 {
		return errors.New("user found")
	}
	return nil
}

// LogRequest logs a *http.Request using the Logger
func LogRequest(r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		Logger.Println(err)
	}
	Logger.Println(fmt.Sprintf("%q", requestDump))
}

// UserToToken convert the user information to a valid token
func UserToToken(u structures.User) string {
	// TODO
	return "token"
}
