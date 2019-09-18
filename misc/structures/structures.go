package structures

import (
	"time"
)

// User struct used to store information in the datastore
type User struct {
	Email     string    `firestore:"email"`
	Pass      string    `firestore:"Hashed_password"`
	Created   time.Time `firestore:"created"`
	Validated bool      `firestore:"valid_user"`
	id        int64     // The integer ID used in the firestore.
}

// Login struct used to parse /login information
type Login struct {
	UserMail string `json:"user_email"`
	UserPass string `json:"user_password"`
}

// Token struct used to parse token information
type Token struct {
	Token string `json:"token"`
}

// Email struct used to parse /verify/email
type Email struct {
	UserMail string `json:"user_email"`
}
