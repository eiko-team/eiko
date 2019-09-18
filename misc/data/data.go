// +build !mock

package data

import (
	"context"
	"errors"
	"log"
	"os"

	"eiko/misc/structures"

	"cloud.google.com/go/datastore"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "data: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

type Data struct {
	// Client is used to take advantage of the datastore api
	Client *datastore.Client
	// Ctx is the context of the Datastore
	Ctx context.Context
	// Users users name inside the datastore
	Users string
}

// InitData return an initialised Data struct
func InitData(projID string) Data {
	var d Data
	d.Users = "Users"
	d.Ctx = context.Background()

	var err error
	d.Client, err = datastore.NewClient(d.Ctx, projID)
	if err != nil {
		Logger.Fatalf("Could not create datastore client: %v", err)
	}

	return d
}

// GetUser is used to find if a email is already used in the datastore
func (d Data) GetUser(UserMail string) (structures.User, error) {
	var users []structures.User
	q := datastore.NewQuery(d.Users).Filter("Email =", UserMail).Limit(1)
	if _, err := d.Client.GetAll(d.Ctx, q, &users); err != nil {
		return structures.User{}, errors.New("Could no fetch users")
	}
	if len(users) == 0 {
		return structures.User{}, errors.New("no user found")
	}
	return users[0], nil
}

// StoreUser is used to store a user in the datastore
func (d Data) StoreUser(user structures.User) error {
	key := datastore.IncompleteKey(d.Users, nil)
	_, err := d.Client.Put(d.Ctx, key, &user)
	return err
}
