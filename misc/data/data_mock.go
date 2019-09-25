// +build mock

package data

import (
	"fmt"
	"time"

	"eiko/misc/hash"
	"eiko/misc/structures"
)

var (
	UserStored      bool
	Logged          bool
	GetUser         bool
	Inited          bool
	GetStore        bool
	StoreStore      bool
	StoreConsumable bool
	GetConsumables  bool
	Error           error
	pass, _         = hash.Hash("pass")
	ErrTest         = fmt.Errorf("Test %s", "error")
	User            = structures.User{}
	UserTest        = structures.User{
		Email:     "test@test.ts",
		Pass:      pass, // hashed password 'pass'
		Created:   time.Now(),
		Validated: false,
	}
	Store     = structures.Store{}
	StoreTest = structures.Store{
		Name:       "test store",
		Address:    "test store",
		Country:    "test store",
		Zip:        "test store",
		UserRating: 5,
	}
	Consumable     = structures.Consumable{}
	ConsumableTest = structures.Consumable{
		Name:    "",
		Company: "",
		Characteristics: structures.Characteristics{
			GlobalInterest: structures.GlobalInterest{
				Boycott:          false,
				EcologicalImpact: "",
				SocialImpact:     "",
			},
			Health: structures.Health{
				Additive: []string{"E404"},
				Allergen: []string{"glutten"},
				Nutrition: structures.Nutrition{
					Energie:       0,
					Fat:           0,
					Fibres:        0,
					Glucides:      0,
					Lipides:       0,
					Proteins:      0,
					Salt:          0,
					SaturatedFat:  0,
					SugarGlucides: 0,
				},
			},
		},
		Pictures: structures.Pictures{
			Back:        "",
			BarCode:     "",
			Composition: "",
			Front:       "",
		},
		Quantity: structures.Quantity{
			Kg:    0,
			Litre: 0,
		},
	}
	Consumables     = structures.Consumables{}
	ConsumablesTest = structures.Consumables{
		Consumable: ConsumableTest,
		Store:      StoreTest,
		Stock:      structures.Stock{},
	}
)

// Data container for all data relative variables
type Data struct {

	// User the user making the request. Got from the cookie in the header
	User structures.User
}

// InitData return an initialised Data struct
func InitData(projID string) Data {
	Inited = true
	var d Data
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

// GetStore is used to find if a store is already in the datastore using
// it's name and location
func (d Data) GetStore(structures.Store) (structures.Store, error) {
	GetStore = true
	return Store, Error
}

// StoreStore is used to store a log in the datastore
func (d Data) StoreStore(store structures.Store) error {
	StoreStore = true
	return Error
}

// StoreConsumable is used to store a log in the datastore
func (d Data) StoreConsumable(consumable structures.Consumable) error {
	StoreConsumable = true
	return Error
}

// GetConsumables is used to store a log in the datastore
func (d Data) GetConsumables(consumable structures.Query) ([]structures.Consumables, error) {
	GetConsumables = true
	return []structures.Consumables{}, Error
}
