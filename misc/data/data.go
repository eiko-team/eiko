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

// Data container for all data relative variables
// initiated inside InitData(...)
type Data struct {
	// client is used to take advantage of the datastore api
	client *datastore.Client
	// ctx is the context of the Datastore
	ctx context.Context
	// users users name inside the datastore
	users string
	// logs log name inside the datastore
	logs string
	// stores stores name inside the datastore
	stores string
	// consumables log name inside the datastore
	consumables string
	// stocks log name inside the datastore
	stocks string

	// User the user making the request. Got from the cookie in the header
	User structures.User
}

// InitData return an initialised Data struct
func InitData(projID string) Data {
	var d Data
	d.users = "Users"
	d.logs = "Logs"
	d.stores = "Stores"
	d.consumables = "Consumables"
	d.stocks = "Stocks"
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

// GetStore is used to find if a store is already in the datastore using
// it's name and location
func (d Data) GetStore(store structures.Store) (structures.Store, error) {
	var stores []structures.Store
	q := datastore.NewQuery(d.stores).
		Filter("Name =", store.Name).
		Filter("Address =", store.Address).
		Filter("Country =", store.Country).
		Filter("Zip =", store.Zip).
		Limit(1)
	if _, err := d.client.GetAll(d.ctx, q, &stores); err != nil {
		return structures.Store{}, errors.New("Could no fetch stores")
	}
	if len(stores) == 0 {
		return structures.Store{}, errors.New("no store found")
	}
	return stores[0], nil
}

// StoreStore is used to store a log in the datastore
func (d Data) StoreStore(store structures.Store) error {
	key := datastore.IncompleteKey(d.stores, nil)
	_, err := d.client.Put(d.ctx, key, &store)
	return err
}

// StoreConsumable is used to store a log in the datastore
func (d Data) StoreConsumable(consumable structures.Consumable) error {
	key := datastore.IncompleteKey(d.consumables, nil)
	_, err := d.client.Put(d.ctx, key, &consumable)
	return err
}
