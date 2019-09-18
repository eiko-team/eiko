package umanagement

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"eiko/misc/data"
	"eiko/misc/hash"
	"eiko/misc/misc"
	"eiko/misc/structures"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "umanagement: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	// Salt salt for the passwrd hashing
	Salt = hash.GenerateKey(666)
)

// Login get the Login informations and return the token to the user if the
// credentials are valid
func Login(d data.Data, r *http.Request) (string, error) {
	var i structures.Login
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("1.0.0")
	}

	user, err := d.GetUser(i.UserMail)
	if err != nil {
		return "", errors.New("1.0.1")
	}

	err = hash.CompareHash(user.Pass, i.UserPass, Salt)
	if err != nil {
		return "", errors.New("1.0.2")
	}

	token, err := misc.UserToToken(user)
	if err != nil {
		return "", errors.New("1.0.3")
	}

	return fmt.Sprintf("{\"token\":\"%s\"}", token), nil
}

// Register adds a new user to the datastore if the credentials are valid
func Register(d data.Data, r *http.Request) (string, error) {
	var i structures.Login
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("1.1.0")
	}

	if _, err = d.GetUser(i.UserMail); err == nil {
		return "", errors.New("1.1.1")
	}

	pass, err := hash.Hash(i.UserPass, Salt)
	if err != nil {
		return "", errors.New("1.1.2")
	}

	user := structures.User{
		Email:     i.UserMail,
		Pass:      pass,
		Created:   time.Now(),
		Validated: false,
	}
	if d.StoreUser(user); err != nil {
		Logger.Println(err)
		return "", errors.New("1.1.3")
	}

	token, err := misc.UserToToken(user)
	if err != nil {
		return "", errors.New("1.1.4")
	}

	return fmt.Sprintf("{\"token\":\"%s\"}", token), nil
}

// Delete delete an account
func Delete(d data.Data, r *http.Request) (string, error) {
	var i structures.Token
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("1.2.0")
	}

	if !misc.ValidateToken(i.Token) {
		return "", errors.New("1.2.1")
	}

	// TODO: acctually delete user

	return "{\"done\":\"true\"}", nil
}

// UpdateToken delete an account
func UpdateToken(d data.Data, r *http.Request) (string, error) {
	var i structures.Token
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("1.3.0")
	}

	if !misc.ValidateToken(i.Token) {
		return "", errors.New("1.3.1")
	}

	user, err := misc.TokenToUser(i.Token)
	if err != nil {
		return "", errors.New("1.3.2")
	}

	token, err := misc.UserToToken(user)
	if err != nil {
		return "", errors.New("1.3.3")
	}
	return fmt.Sprintf("{\"token\":\"%s\"}", token), nil
}
