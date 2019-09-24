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

// Log struct used to store Logs in the datastore
type Log struct {
	Email   string    `firestore:"email"`
	Log     string    `firestore:"log"`
	Created time.Time `firestore:"created"`
	id      int64     // The integer ID used in the firestore.
}

// Logging struct used to parse /log information
type Logging struct {
	Log   string `json:"message"`
	Token string `json:"user_token"`
}

// Login struct used to parse /login information
type Login struct {
	UserMail string `json:"user_email"`
	UserPass string `json:"user_password"`
}

// Email struct used to parse /verify/email
type Email struct {
	UserMail string `json:"user_email"`
}

// Password struct used to parse /verify/password
type Password struct {
	Password string `json:"password"`
}

// Store struct used to store store data
type Store struct {
	Name       string `json:"name" firestore:"Name"`
	Address    string `json:"address" firestore:"Address"`
	Country    string `json:"country" firestore:"Country"`
	Zip        string `json:"zip" firestore:"Zip"`
	UserRating int    `json:"user_rating" firestore:"created"`
	id         int64  // The integer ID used in the firestore.
}

// Consumable struct used to parse /consumable/...
type Consumable struct {
	Name            string `json:"name"`
	Company         string `json:"Compagny"`
	Characteristics struct {
		GlobalInterest struct {
			Boycott          bool   `json:"boycott"`
			EcologicalImpact string `json:"ecological_impact"`
			SocialImpact     string `json:"social_impact"`
		} `json:"global_interest"`
		Health struct {
			Additive  []string `json:"Additive"`
			Allergen  []string `json:"allergen"`
			Nutrition struct {
				Energie       float64 `json:"energie"`
				Fat           float64 `json:"fat"`
				Fibres        float64 `json:"fibres"`
				Glucides      float64 `json:"glucides"`
				Lipides       float64 `json:"lipides"`
				Proteines     float64 `json:"proteines"`
				Salt          float64 `json:"salt"`
				SaturatedFat  float64 `json:"saturated_fat"`
				SugarGlucides float64 `json:"sugar_glucides"`
			} `json:"nutrition"`
		} `json:"health"`
	} `json:"characteristics"`
	Pictures struct {
		Back        string `json:"back"`
		BarCode     string `json:"bar_code"`
		Composition string `json:"composition"`
		Front       string `json:"front"`
	} `json:"pictures"`
	Quantity struct {
		Kg    int `json:"kg"`
		Litre int `json:"litre"`
	} `json:"quantity"`
}

// Stock TODO
type Stock struct {
}

// ConsumableWithEverything struct used to parse /consumable/...
type ConsumableWithEverything struct {
	Consumable Consumable `json:"consumable"`
	Store      Store      `json:"store"`
	Stock      Stock      `json:"Stock"`
}
