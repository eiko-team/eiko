// +build mock

package data

import (
	"fmt"
	"time"

	"eiko/misc/structures"
)

var (
	UserStored bool
	Logged     bool
	GetUser    bool
	Inited     bool
	TestError  = fmt.Errorf("Test %s", "error")
	Error      error
	User       = structures.User{}
	TestUser   = structures.User{
		Email:     "test@test.ts",
		Pass:      "$2a$10$EVCZ/75E1TCgpOZFypJC4ejYDDTPk9lAGwLKGhp6jESMWfl/4Bl/e", // hashed password 'pass'
		Created:   time.Now(),
		Validated: false,
	}
)

type Data struct {
	// Users users name inside the datastore
	Users string
}

// InitData return an initialised Data struct
func InitData(projID string) Data {
	Inited = true
	var d Data
	d.Users = "Users"
	return d
}

// GetUser is used to find if a email is already used in the datastore
func (d Data) GetUser(UserMail string) (structures.User, error) {
	GetUser = true
	return User, Error
}

// StoreUser is used to store a user in the datastore
func (d Data) StoreUser(user structures.User) error {
	UserStored = true
	return nil
}

// Log is used to store a log in the datastore
func (d Data) Log(user structures.Log) error {
	Logged = true
	return nil
}
