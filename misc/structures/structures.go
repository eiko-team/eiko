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

// GlobalInterest is the global interest of a consumable
type GlobalInterest struct {
	Boycott          bool   `json:"boycott"`
	EcologicalImpact string `json:"ecological_impact"`
	SocialImpact     string `json:"social_impact"`
}

// Nutrition Nutrition facts on the consumable
type Nutrition struct {
	Energie       float64 `json:"energie"`
	Fat           float64 `json:"fat"`
	Fibres        float64 `json:"fibres"`
	Glucides      float64 `json:"glucides"`
	Lipides       float64 `json:"lipides"`
	Proteins      float64 `json:"proteins"`
	Salt          float64 `json:"salt"`
	SaturatedFat  float64 `json:"saturated_fat"`
	SugarGlucides float64 `json:"sugar_glucides"`
}

// Health Health status of a consumable
type Health struct {
	Additive  []string  `json:"Additive"`
	Allergen  []string  `json:"allergen"`
	Nutrition Nutrition `json:"nutrition"`
}

// Characteristics Characteristics of a consumable
type Characteristics struct {
	GlobalInterest GlobalInterest `json:"global_interest"`
	Health         Health         `json:"health"`
}

// Pictures all Pictures needed for a consumable
type Pictures struct {
	Back        string `json:"back"`
	BarCode     string `json:"bar_code"`
	Composition string `json:"composition"`
	Front       string `json:"front"`
}

// Quantity Quantity of the product in a pack
type Quantity struct {
	Kg    int `json:"kg"`
	Litre int `json:"litre"`
}

// Consumable struct used to parse /consumable/...
type Consumable struct {
	Name            string          `json:"name"`
	Company         string          `json:"Compagny"`
	Characteristics Characteristics `json:"characteristics"`
	Pictures        Pictures        `json:"pictures"`
	Quantity        Quantity        `json:"quantity"`
}

// Stock Stock of a product in a store
type Stock struct {
	PackQuantity int `json:"pack_quantity" firestore:"pack_quantity"`
	NbPacks      int `json:"nb_packs" firestore:"nb_packs"`
	PackPrice    int `json:"pack_price" firestore:"pack_price"`
	Available    int `json:"available" firestore:"available"`
}

// Consumables struct used to parse /consumable/...
type Consumables struct {
	Consumable Consumable `json:"consumable"`
	Store      Store      `json:"store"`
	Stock      Stock      `json:"Stock"`
}

// Location struct to store location informations
// https://www.w3schools.com/html/html5_geolocation.asp
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Query used to query certain api to get a personalized result
type Query struct {
	Query    string   `json:"query"`
	Limit    int      `json:"limit"`
	Location Location `json:"limit"`
}
