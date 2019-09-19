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
	// client is used to take advantage of the datastore api
	client *datastore.Client
	// ctx is the context of the Datastore
	ctx context.Context
	// users users name inside the datastore
	users string
	// logs log name inside the datastore
	logs string
}

// InitData return an initialised Data struct
func InitData(projID string) Data {
	var d Data
	d.users = "Users"
	d.logs = "Logs"
	d.ctx = context.Background()

	var err error
	d.client, err = datastore.NewClient(d.ctx, projID)
	if err != nil {
		Logger.Fatalf("Could not create datastore client: %v", err)
	}

	return d
}

// GetUser is used to find if a email is already used in the datastore
func (d Data) GetUser(UserMail string) (structures.User, error) {
	var users []structures.User
	q := datastore.NewQuery(d.users).Filter("Email =", UserMail).Limit(1)
	if _, err := d.client.GetAll(d.ctx, q, &users); err != nil {
		return structures.User{}, errors.New("Could no fetch users")
	}
	if len(users) == 0 {
		return structures.User{}, errors.New("no user found")
	}
	return users[0], nil
}

// StoreUser is used to store a user in the datastore
func (d Data) StoreUser(user structures.User) error {
	key := datastore.IncompleteKey(d.users, nil)
	_, err := d.client.Put(d.ctx, key, &user)
	return err
}

// Log is used to store a log in the datastore
func (d Data) Log(log structures.Log) error {
	key := datastore.IncompleteKey(d.logs, nil)
	_, err := d.client.Put(d.ctx, key, &log)
	return err
}
