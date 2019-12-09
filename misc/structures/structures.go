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

// MergeUser Merge two user into one:
// i2 fill empty fields of i1
func MergeUser(i1, i2 User) User {
	if i1.Email == "" {
		i1.Email = i2.Email
	}
	if i1.Pass == "" {
		i1.Pass = i2.Pass
	}
	if i1.Created == time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) {
		i1.Created = i2.Created
	}
	if i1.Validated == false {
		i1.Validated = i2.Validated
	}
	if i1.ID == 0 {
		i1.ID = i2.ID
	}
	return i1
}

// Log struct used to store Logs in the datastore
type Log struct {
	Email   string    `firestore:"email"`
	Log     string    `firestore:"log"`
	Created time.Time `firestore:"created"`
	IP      string    `firestore:"ip"`
	Port    string    `firestore:"port"`
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
	// score averall score of the store, used to give a feeback score
	Score int `json:"score" firestore:"score"`
	// scoreNb number of people scoring the store
	ScoreNb int   `json:"score_nb" firestore:"scoreNb"`
	ID      int64 // The integer ID used in the firestore.
}

// Consumable struct used to parse /consumable/...
type Consumable struct {
	// The integer ID used in the firestore.
	ID            int64
	Name          string `json:"name"`
	Company       string `json:"compagny"`
	Manufacturing string `json:"manifacturing"`
	Created       time.Time
	Creator       int64
	NewVersion    int64
	Source        string
	Code          []string `json:"code"`
	Categories    []string `json:"categories"`
	Tags          []string `json:"tags"`
	Packaging     []string `json:"packaging"`
	// Nutrition facts on the consumable
	Fat           float64 `json:"fat"`
	Fiber         float64 `json:"fiber"`
	Glucides      float64 `json:"glucides"`
	Proteins      float64 `json:"proteins"`
	Sodium        float64 `json:"sodium"`
	SaturatedFat  float64 `json:"saturated_fat"`
	SugarGlucides float64 `json:"sugar_glucides"`
	Energy        float64 `json:"energy"`
	Alcool        float64 `json:"alcool"`
	// Health status of the consumable
	Additive   []string `json:"additive"`
	Ingredient []string `json:"ingredient"`
	Vitamins   []string `json:"vitamins"`
	Allergen   []string `json:"allergen"`
	NutriScore string   `json:"nutri_score"`
	// all Pictures needed for the consumable
	Back        string `json:"back"`
	Composition string `json:"composition"`
	Front       string `json:"front"`
	// Quantity of the product in a pack
	Grammes int `json:"grammes"`
	MLitre  int `json:"mililitre"`
	// Quality of the product
	Label []string `json:"label"`
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

// ConsumablesID struct used to get a Consumable
type ConsumablesID struct {
	ConsumableID int64 `json:"consumable"`
	StoreID      int64 `json:"store"`
	StockID      int64 `json:"stock"`
}

// ListContent list content
type ListContent struct {
	ID          int64         // The integer ID used in the firestore.
	ListID      int64         `json:"list_id" firestore:"list_id"`
	Consumables ConsumablesID `json:"consumable" firestore:"consumable"`
	Name        string        `json:"name" firestore:"name"`
	Done        bool          `json:"done" firestore:"done"`
	Erased      bool          `json:"erased" firestore:"erased"`
	// Mode is used to signify the mode of the consumable
	// It's content should be:
	// 	- "sample": for testing purpose
	// 	- "consumable": for a "real" consumable
	// 	- "personnal": for a user imported consmable where only the Name field is
	// used. stock and store are also not used
	Mode string `json:"mode"`
}
