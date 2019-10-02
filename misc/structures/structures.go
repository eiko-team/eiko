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
	ID        int64     // The integer ID used in the firestore.
}

// Log struct used to store Logs in the datastore
type Log struct {
	Email   string    `firestore:"email"`
	Log     string    `firestore:"log"`
	Created time.Time `firestore:"created"`
	ID      int64     // The integer ID used in the firestore.
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
	GeoHash    int    `json:"geohash" firestore:"geohash"`
	ID         int64  // The integer ID used in the firestore.
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
	ID              int64           // The integer ID used in the firestore.
	// Mode is used to signifiy the mode of the consumable
	// It's content should be:
	// 	- "sample": for testing purpose
	// 	- "consumable": for a "real" consumable
	// 	- "personnal": for a user imported consmable where only the Name field is
	// used. stock and store are also not used
	Mode string `json:"mode"`
}

// Stock Stock of a product in a store
type Stock struct {
	ID            int64 // The integer ID used in the firestore.
	PackQuantity  int   `json:"pack_quantity" firestore:"pack_quantity"`
	NbPacks       int   `json:"nb_packs" firestore:"nb_packs"`
	PackPrice     int   `json:"pack_price" firestore:"pack_price"`
	Available     bool  `json:"available" firestore:"available"`
	StoreKey      int   `json:"store_key" firestore:"store_key"`
	ConsumableKey int   `json:"consumable_key" firestore:"consumable_key"`
	GeoHash       int   `json:"geohash" firestore:"geohash"`
}

// Consumables struct used to parse /consumable/...
type Consumables struct {
	Consumable Consumable `json:"consumable"`
	Store      Store      `json:"store"`
	Stock      Stock      `json:"stock"`
}

// Query used to query certain api to get a personalized result
type Query struct {
	Query     string  `json:"query"`
	Limit     int     `json:"limit"`
	Size      uint64  `json:"size"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// List content of a list
type List struct {
	ID   int64  `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
}

// ListOwner list reference to know who own the list
type ListOwner struct {
	ID     int64  // The integer ID used in the firestore.
	ListID int64  `json:"list_id" firestore:"list_id"`
	Email  string `json:"email" firestore:"email"`
}

// ListContent list content
type ListContent struct {
	ID          int64       // The integer ID used in the firestore.
	ListID      int64       `json:"list_id" firestore:"list_id"`
	Consumables Consumables `json:"consumable" firestore:"consumable"`
	Name        string      `json:"name" firestore:"name"`
	Done        bool        `json:"done" firestore:"done"`
	Erased      bool        `json:"erased" firestore:"erased"`
}
