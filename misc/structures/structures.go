package structures

import (
	"time"
)

type User struct {
	Email     string    `firestore:"email"`
	Pass      string    `firestore:"Hashed_password"`
	Created   time.Time `firestore:"created"`
	Validated bool      `firestore:"valid_user"`
	id        int64     // The integer ID used in the firestore.
}

type Login struct {
	UserMail string `json:"user_email"`
	UserPass string `json:"user_password"`
}
