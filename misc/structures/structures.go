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
	if !i1.Validated {
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

// MergeStore Merge two user into one:
// i2 fill empty fields of i1
func MergeStore(i1, i2 Store) Store {
	if i1.Country == "" {
		i1.Country = i2.Country
	}
	i1.Address = i2.Address
	i1.Name = i2.Name
	i1.Zip = i2.Zip
	i1.GeoHash = i2.GeoHash
	return i1
}

// MergeStrings put the content of s2 in s1 if s1 is empty
func MergeStrings(s1, s2 []string) []string {
	if len(s1) == 0 {
		return append(s1, s2...)
	}
	return s1
}

// MergeInt put the content of i2 in i1 if i1 is empty
func MergeInt(i1, i2 int) int {
	if i1 == 0 {
		return i2
	}
	return i1
}

// MergeFloat put the content of f2 in f1 if f1 is empty
func MergeFloat(f1, f2 float64) float64 {
	if f1 == 0 {
		return f2
	}
	return f1
}

// MergeString put the content of i2 in i1 if i1 is empty
func MergeString(s1, s2 string) string {
	if s1 == "" {
		return s2
	}
	return s1
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

// MergeConsumable Merge two consumable into one:
// C2 fill empty fields of C1
func MergeConsumable(c1, c2 Consumable) Consumable {
	c1.Additive = MergeStrings(c1.Additive, c2.Additive)
	c1.Alcool = MergeFloat(c1.Alcool, c2.Alcool)
	c1.Allergen = MergeStrings(c1.Allergen, c2.Allergen)
	c1.Back = MergeString(c1.Back, c2.Back)
	c1.Categories = MergeStrings(c1.Categories, c2.Categories)
	c1.Code = MergeStrings(c1.Code, c2.Code)
	c1.Company = MergeString(c1.Company, c2.Company)
	c1.Composition = MergeString(c1.Composition, c2.Composition)
	c1.Energy = MergeFloat(c1.Energy, c2.Energy)
	c1.Fat = MergeFloat(c1.Fat, c2.Fat)
	c1.Fiber = MergeFloat(c1.Fiber, c2.Fiber)
	c1.Front = MergeString(c1.Front, c2.Front)
	c1.Glucides = MergeFloat(c1.Glucides, c2.Glucides)
	c1.Ingredient = MergeStrings(c1.Ingredient, c2.Ingredient)
	c1.Label = MergeStrings(c1.Label, c2.Label)
	c1.Manufacturing = MergeString(c1.Manufacturing, c2.Manufacturing)
	c1.NutriScore = MergeString(c1.NutriScore, c2.NutriScore)
	c1.Packaging = MergeStrings(c1.Packaging, c2.Packaging)
	c1.Proteins = MergeFloat(c1.Proteins, c2.Proteins)
	c1.SaturatedFat = MergeFloat(c1.SaturatedFat, c2.SaturatedFat)
	c1.Sodium = MergeFloat(c1.Sodium, c2.Sodium)
	c1.Source = MergeString(c1.Source, c2.Source)
	c1.SugarGlucides = MergeFloat(c1.SugarGlucides, c2.SugarGlucides)
	c1.Tags = MergeStrings(c1.Tags, c2.Tags)
	c1.Vitamins = MergeStrings(c1.Vitamins, c2.Vitamins)
	return c1
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
