// +build mock

package data

import (
	"fmt"
	"time"

	"github.com/eiko-team/eiko/misc/hash"
	"github.com/eiko-team/eiko/misc/structures"
)

var (
	StoreUser       bool
	Log             bool
	GetUser         bool
	Inited          bool
	GetStore        bool
	StoreStore      bool
	StoreConsumable bool
	GetConsumables  bool
	GetList         bool
	CreateList      bool
	GetAllLists     bool
	GetListContent  bool
	StoreContent    bool
	Error           error
	Error2          error
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
	StoreRe   = `{"name":"[a-z ]+","address":"[a-z ]+","country":"[a-z ]+","zip":"[a-z ]+","user_rating":\d+,"geohash":\d+,"ID":\d+}`
	StoreTest = structures.Store{
		Name:       "test store",
		Address:    "test store",
		Country:    "test store",
		Zip:        "test store",
		UserRating: 5,
	}
	Stock          = structures.Stock{}
	StockRe        = `{"ID":\d+,"pack_quantity":\d+,"nb_packs":\d+,"pack_price":\d+,"available":[a-z]+,"store_key":\d+,"consumable_key":\d+,"geohash":\d+}`
	StockTest      = structures.Stock{}
	Consumable     = structures.Consumable{}
	ConsumableRe   = `{"ID":\d+,"name":"[a-zA-Z 0-9,_]+","compagny":"[a-zA-Z 0-9,_]+","manifacturing":"[a-zA-Z 0-9,_]+","Created":"0001-01-01T00:00:00Z","Creator":\d+,"NewVersion":\d+,"Source":"[a-zA-Z 0-9,_]*","code":\[("[a-zA-Z 0-9,_-]+",?)*\],"categories":\[("[a-zA-Z 0-9,_-]+",?)*\],"tags":\[("[a-zA-Z 0-9,_-]+",?)*\],"packaging":\[("[a-zA-Z 0-9,_-]+",?)*\],"fat":\d+,"fiber":\d+,"glucides":\d+,"proteins":\d+,"sodium":\d+,"saturated_fat":\d+,"sugar_glucides":\d+,"energy":\d+,"alcool":\d+,"additive":\[("[a-zA-Z 0-9,_-]+",?)*\],"ingredient":\[("[a-zA-Z 0-9,_-]+",?)*\],"vitamins":\[("[a-zA-Z 0-9,_-]+",?)*\],"allergen":\[("[a-zA-Z 0-9,_-]+",?)*\],"nutri_score":"[a-z]","back":"url","composition":"url","front":"url","grammes":\d+,"mililitre":\d+,"label":\[("[a-zA-Z 0-9,_-]+",?)*\]}`
	ConsumableTest = structures.Consumable{
		Name:          "Simple Name",
		Company:       "Simple Compagny Name",
		Additive:      []string{"E404"},
		Allergen:      []string{"glutten"},
		Fat:           9001,
		Fiber:         9001,
		Glucides:      9001,
		Proteins:      9001,
		Sodium:        9001,
		SaturatedFat:  9001,
		SugarGlucides: 9001,
		Energy:        9001,
		Alcool:        9001,
		Back:          "url",
		Composition:   "url",
		Front:         "url",
		Manufacturing: "1 rue du pont, 75001 Paris",
		Code:          []string{"code-13", "0001006530676", "000100653067x", "00010065306xx", "0001006530xxx", "000100653xxxx", "00010065xxxxx", "0001006xxxxxx", "000100xxxxxxx", "00010xxxxxxxx", "0001xxxxxxxxx", "000xxxxxxxxxx", "00xxxxxxxxxxx", "0xxxxxxxxxxxx"},
		Categories:    []string{"Conserves"},
		Tags:          []string{"paste", "vegetarien", "terre", "soybean"},
		Packaging:     []string{},
		Ingredient:    []string{"WATER", "SOYBEANS", "RICE", "SALT", "ALCOHOL"},
		Vitamins:      []string{"B12"},
		NutriScore:    "a",
		Grammes:       0,
		MLitre:        1000,
		Label:         []string{"AOC"},
	}
	Consumables   = []structures.Consumables{}
	ConsumablesRe = fmt.Sprintf("{\"consumable\":%s,\"store\":%s,\"stock\":%s}",
		ConsumableRe, StoreRe, StockRe)
	ConsumablesTest = []structures.Consumables{
		{
			Consumable: ConsumableTest,
			Store:      StoreTest,
			Stock:      structures.Stock{},
		},
	}
	List     = structures.List{}
	ListRe   = `{"id":\d+,"name":"[a-zA-Z0-9 ]+"}`
	ListTest = structures.List{
		ID:   0,
		Name: "List name Test",
	}
	ConsumableIDsRe = `{"consumable":0,"store":0,"stock":0}`
	ListContent     = structures.ListContent{}
	ListContentRe   = fmt.Sprintf(`{"ID":\d+,"list_id":\d+,"consumable":%s,"name":"[a-zA-Z0-9 ]+","done":false,"erased":false,"mode":"[a-z]+"}`,
		ConsumableIDsRe)
	ListContentTest = structures.ListContent{
		Consumables: structures.ConsumablesID{
			ConsumableID: 12,
			StoreID:      24,
			StockID:      48,
		},
		Name: ConsumablesTest[0].Consumable.Name,
		Mode: "sample",
	}
	ID     = int64(0)
	IDTest = int64(42)
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
	StoreUser = true
	return Error2
}

// Log is used to store a log in the datastore
func (d Data) Log(user structures.Log) error {
	Log = true
	return Error
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
func (d Data) StoreConsumable(consumable structures.Consumable) (int64, error) {
	StoreConsumable = true
	return ID, Error
}

// GetConsumables is used to store a log in the datastore
func (d Data) GetConsumables(consumable structures.Query) ([]structures.Consumables, error) {
	GetConsumables = true
	return Consumables, Error
}

func (d Data) GetList(id int64) (structures.List, error) {
	GetList = true
	return List, Error
}

func (d Data) CreateList(list structures.List) (structures.List, error) {
	CreateList = true
	return List, Error
}

func (d Data) GetAllLists() ([]structures.List, error) {
	GetAllLists = true
	return []structures.List{List}, Error
}

// GetListContent return list content
func (d Data) GetListContent(id int64) ([]structures.ListContent, error) {
	GetListContent = true
	return []structures.ListContent{ListContent}, Error
}

// StoreContent is used to store a content in the datastore
func (d Data) StoreContent(content structures.ListContent) (int64, error) {
	StoreContent = true
	return ID, Error
}
