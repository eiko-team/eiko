package user_managment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"eiko/misc/hash"
	"eiko/misc/misc"
	"eiko/misc/structures"

	"cloud.google.com/go/datastore"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "user_managment: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	// Salt salt for the passwrd hashing
	Salt = "tmp"

	// Users users name inside the datastore
	Users = "Users"
)

func ParseJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(v)
	if err != nil {
		Logger.Printf("\033[31mError\033[0m: '%s'\n", err.Error())
	}
	return err
}

func Login(w http.ResponseWriter, r *http.Request, ctx context.Context,
	client *datastore.Client) (string, error) {
	var i structures.Login
	err := ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("1.0.0")
	}
	var users []structures.User
	q := datastore.NewQuery(Users).Filter("Email =", i.UserMail).Limit(1)
	if _, err := client.GetAll(ctx, q, &users); err != nil {
		Logger.Printf("%q", err)
		return "", errors.New("1.0.1")
	}
	if len(users) == 0 {
		Logger.Println("no user found")
		return "", errors.New("1.0.2")
	}
	User := users[0]
	Logger.Printf("%v\n", User)

	err = hash.CompareHash(User.Pass, i.UserPass, Salt)
	if err != nil {
		return "", errors.New("1.0.3")
	}
	return fmt.Sprintf("{\"token\":\"%s\"}", misc.UserToToken(User)), nil
}

func Register(w http.ResponseWriter, r *http.Request, ctx context.Context,
	client *datastore.Client) (string, error) {
	var i structures.Login
	err := ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("1.1.0")
	}

	if err = misc.UniqEmail(Users, i.UserMail, ctx, client); err != nil {
		return "", errors.New("1.1.1\", \"message\":\"not a unique email")
	}

	pass, err := hash.Hash(i.UserPass, Salt)
	if err != nil {
		return "", errors.New("1.1.2")
	}

	User := structures.User{
		Email:     i.UserMail,
		Pass:      pass,
		Created:   time.Now(),
		Validated: false,
	}
	key := datastore.IncompleteKey(Users, nil)
	_, err = client.Put(ctx, key, &User)
	if err != nil {
		Logger.Println(err)
		return "", errors.New("1.1.3")
	}
	return fmt.Sprintf("{\"token\":\"%s\"}", misc.UserToToken(User)), nil
}