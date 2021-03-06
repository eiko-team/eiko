// +build !mock

package data

import (
	"context"
	"errors"
	"os"
	"reflect"

	"github.com/eiko-team/eiko/misc/log"
	"github.com/eiko-team/eiko/misc/math"
	"github.com/eiko-team/eiko/misc/research"
	"github.com/eiko-team/eiko/misc/structures"
	"github.com/eiko-team/eiko/misc/user"

	"cloud.google.com/go/datastore"
	// https://blog.nobugware.com/post/2015/leveldb_geohash_golang/
	"github.com/mmcloughlin/geohash"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "data: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	// GeoHashRadius is the query bounding box max value
	GeoHashRadius = uint64(500)
)

// Data container for all data relative variables
// initiated inside InitData(...)
type Data struct {
	// client is used to take advantage of the datastore api
	client *datastore.Client
	// ctx is the context of the Datastore
	ctx context.Context
	// users users name inside the datastore
	users string
	// logs log name inside the datastore
	logs string
	// stores stores name inside the datastore
	stores string
	// consumables log name inside the datastore
	consumables string
	// stocks log name inside the datastore
	stocks string
	// listsOwner list of lists and owners
	listsOwner string
	// list list of lists and owners
	list string
	// listContent list content
	listContent string

	// User the user making the request. Got from the cookie in the header
	User structures.User

	// rs search engine
	rs research.Research
}

// InitData return an initialised Data struct
func InitData(projID string) Data {
	var d Data
	d.users = "Users"
	d.logs = "Logs"
	d.stores = "Stores"
	d.consumables = "Consumables"
	d.stocks = "Stocks"
	d.listsOwner = "listsOwner"
	d.list = "list"
	d.listContent = "ListContent"
	d.ctx = context.Background()

	var err error
	d.client, err = datastore.NewClient(d.ctx, projID)
	if err != nil {
		Logger.Fatalf("Could not create datastore client: %v", err)
	}

	d.rs = research.InitReseach()

	return d
}

// GetUser is used to find if a email is already used in the datastore
func (d Data) GetUser(UserMail string) (structures.User, error) {
	var users []structures.User
	q := datastore.NewQuery(d.users).Filter("Email =", UserMail).Limit(1)
	keys, err := d.client.GetAll(d.ctx, q, &users)
	if err != nil || len(users) == 0 {
		return structures.User{}, errors.New("Could no fetch users")
	}
	users[0].ID = keys[0].ID
	return users[0], nil
}

// StoreUser is used to store a user in the datastore
func (d Data) StoreUser(user structures.User) error {
	user.ID = 0
	key := datastore.IncompleteKey(d.users, nil)
	_, err := d.client.Put(d.ctx, key, &user)
	return err
}

// Log is used to store a log in the datastore
func (d Data) Log(log structures.Log) error {
	key := datastore.IncompleteKey(d.logs, nil)
	_, err := d.client.Put(d.ctx, key, &log)
	return err
}

// GetStore is used to find if a store is already in the datastore using
// it's name and location
func (d Data) GetStore(store structures.Store) (structures.Store, error) {
	var stores []structures.Store
	q := datastore.NewQuery(d.stores).
		Filter("Name =", store.Name).
		Filter("Address =", store.Address).
		Filter("Country =", store.Country).
		Filter("Zip =", store.Zip).
		Limit(1)
	keys, err := d.client.GetAll(d.ctx, q, &stores)
	if err != nil || len(stores) == 0 {
		return structures.Store{}, errors.New("Could no fetch stores")
	}
	stores[0].ID = keys[0].ID
	return stores[0], nil
}

// StoreStore is used to store a consumable in the datastore
func (d Data) StoreStore(store structures.Store) error {
	store.ID = 0
	key := datastore.IncompleteKey(d.stores, nil)
	_, err := d.client.Put(d.ctx, key, &store)
	return err
}

func (d Data) getSameConsumable(consumable structures.Consumable) (structures.Consumable, error) {
	var cons []structures.Consumable
	q := datastore.NewQuery(d.consumables).
		Filter("Name =", consumable.Name).
		Limit(1)
	keys, err := d.client.GetAll(d.ctx, q, &cons)
	if err != nil || len(cons) == 0 {
		return structures.Consumable{}, errors.New("Could no fetch stores")
	}
	cons[0].ID = keys[0].ID
	return cons[0], nil
}

func (d Data) storeConsumable(consumable, prev structures.Consumable) (int64, error) {
	key := datastore.IncompleteKey(d.consumables, nil)
	key, err := d.client.Put(d.ctx, key, &consumable)
	if err != nil {
		return key.ID, err
	}
	consumable.ID = key.ID
	if !reflect.DeepEqual(prev, structures.Consumable{}) {
		prev.NewVersion = consumable.ID
		key := datastore.IDKey(d.consumables, prev.ID, nil)
		d.client.Put(d.ctx, key, &prev)
	}
	err = d.rs.StoreConsumable(consumable)
	return key.ID, err
}

// StoreConsumable is used to store a consumable in the datastore
func (d Data) StoreConsumable(consumable structures.Consumable) (int64, error) {
	cons, err := d.getSameConsumable(consumable)
	if err != nil {
		return d.storeConsumable(consumable, structures.Consumable{})
	}
	consumable = structures.MergeConsumable(cons, consumable)
	return d.storeConsumable(consumable, cons)
}

func (d Data) fetchStock(geo uint64, filter, order string, limit int) ([]structures.Stock, error) {
	var res []structures.Stock
	q := datastore.NewQuery(d.stocks).
		Filter(filter, geo).
		Order(order).
		Limit(limit * 10)
	keys, err := d.client.GetAll(d.ctx, q, &res)
	if err != nil {
		return []structures.Stock{}, errors.New("Could no fetch stocks")
	}
	for i, k := range keys {
		res[i].ID = k.ID
	}
	return res, nil
}

// getOneConsumable return the last version of this consumable
func (d Data) getOneConsumable(id int64) (structures.Consumable, error) {
	var consumable []structures.Consumable
	q := datastore.NewQuery(d.consumables).
		Filter("__key__ =", datastore.IDKey(d.consumables, id, nil)).
		Limit(1)
	keys, err := d.client.GetAll(d.ctx, q, &consumable)
	if err != nil || len(consumable) == 0 {
		return structures.Consumable{}, errors.New("Could no fetch consumable")
	}
	consumable[0].ID = keys[0].ID
	return consumable[0], nil
}

// GetConsumablesTmp is used to find some consumables in the datastore but for
// testing purposes only
func (d Data) GetConsumablesTmp(query structures.Query) ([]structures.Consumables, error) {
	Logger.Println(query)
	IDs, err := d.rs.SearchConsumable(query.Query)
	if err != nil {
		return nil, err
	}

	res := make([]structures.Consumables, len(IDs))
	pos := 0
	for _, id := range IDs {
		c, err := d.getOneConsumable(id)
		if err != nil {
			return nil, err
		}
		if user.IsGood(d.User, c) {
			res[pos] = structures.Consumables{Consumable: c}
			pos++
		}
	}
	return res[:pos], nil
}

// GetConsumables is used to find some consumables in the datastore
func (d Data) GetConsumables(query structures.Query) ([]structures.Consumables, error) {
	geo := geohash.EncodeInt(query.Latitude, query.Longitude)
	stocks1, err := d.fetchStock(geo, "geohash <", "geohash", query.Limit)
	if err != nil {
		return []structures.Consumables{}, err
	}
	stocks2, err := d.fetchStock(geo, "geohash >", "-geohash", query.Limit)
	if err != nil {
		return []structures.Consumables{}, err
	}
	var res []structures.Consumables
	for _, stock := range append(stocks1, stocks2...) {
		var consumable []structures.Consumable
		q := datastore.NewQuery(d.consumables).
			Filter("__key__ =", stock.ConsumableKey).
			Filter("name = ", query.Query).
			Limit(1)
		keys, err := d.client.GetAll(d.ctx, q, &consumable)
		if err != nil || len(consumable) == 0 {
			continue
		}
		consumable[0].ID = keys[0].ID
		var store []structures.Store
		q = datastore.NewQuery(d.stores).
			Filter("__key__ =", stock.StoreKey).
			Limit(1)
		keys, err = d.client.GetAll(d.ctx, q, &store)
		if err != nil || len(store) == 0 {
			continue
		}
		store[0].ID = keys[0].ID
		res = append(res, structures.Consumables{
			Consumable: consumable[0],
			Store:      store[0],
			Stock:      stock,
		})
	}
	return res, nil
}

// GetList return a list if the id provided is right and the user know the list
func (d Data) GetList(id int64) (structures.List, error) {
	// Checks if the user know the list
	var listsOwners []structures.ListOwner
	q := datastore.NewQuery(d.listsOwner).
		Filter("ListID =", id).
		Filter("Email =", d.User.Email).
		Limit(1)
	_, err := d.client.GetAll(d.ctx, q, &listsOwners)
	if err != nil || len(listsOwners) == 0 {
		return structures.List{}, errors.New("No List found")
	}

	// Fetch list
	var lists []structures.List
	k := datastore.IDKey(d.list, id, nil)
	q = datastore.NewQuery(d.list).
		Filter("__key__ =", k).
		Limit(1)
	keys, err := d.client.GetAll(d.ctx, q, &lists)
	if err != nil || len(lists) == 0 {
		return structures.List{}, errors.New("Could no fetch list")
	}
	lists[0].ID = keys[0].ID
	return lists[0], nil
}

// CreateList is used to create a list in the datastore
func (d Data) CreateList(list structures.List) (structures.List, error) {
	list.ID = 0
	key := datastore.IncompleteKey(d.list, nil)
	key, err := d.client.Put(d.ctx, key, &list)
	if err != nil {
		return structures.List{}, err
	}
	list.ID = key.ID
	listOwner := structures.ListOwner{
		ListID: key.ID,
		Email:  d.User.Email,
	}
	key = datastore.IncompleteKey(d.listsOwner, nil)
	_, err = d.client.Put(d.ctx, key, &listOwner)
	return list, err
}

// GetAllLists return all lists of a user
func (d Data) GetAllLists() ([]structures.List, error) {
	var listsOwners []structures.ListOwner
	q := datastore.NewQuery(d.listsOwner).
		Filter("Email =", d.User.Email)
	_, err := d.client.GetAll(d.ctx, q, &listsOwners)
	if err != nil || len(listsOwners) == 0 {
		return []structures.List{}, errors.New("No List found")
	}
	var res []structures.List
	for _, listOwner := range listsOwners {
		var lists []structures.List
		k := datastore.IDKey(d.list, listOwner.ListID, nil)
		q = datastore.NewQuery(d.list).
			Filter("__key__ =", k).
			Limit(1)
		keys, err := d.client.GetAll(d.ctx, q, &lists)
		if err != nil {
			Logger.Println(err)
			return []structures.List{}, nil
		}
		if len(keys) == 0 {
			continue
		}
		res = append(res, structures.List{
			ID:   keys[0].ID,
			Name: lists[0].Name,
		})
	}
	return res, nil
}

// GetListContent return list content
func (d Data) GetListContent(id int64) ([]structures.ListContent, error) {
	// Checks if the user know the list
	var listsOwners []structures.ListOwner
	Logger.Printf("list: %d, user: %+v", id, d.User)
	q := datastore.NewQuery(d.listsOwner).
		Filter("ListID =", id).
		Filter("Email =", d.User.Email).
		Limit(1)
	_, err := d.client.GetAll(d.ctx, q, &listsOwners)
	if err != nil || len(listsOwners) == 0 {
		return []structures.ListContent{}, errors.New("No List found")
	}

	// Fetch content
	var listContents []structures.ListContent
	q = datastore.NewQuery(d.listContent).
		Filter("ListID =", id)
	keys, err := d.client.GetAll(d.ctx, q, &listContents)
	if err != nil || len(listContents) == 0 {
		Logger.Println("No Content found")
		return []structures.ListContent{}, nil
	}
	for i, k := range keys {
		listContents[i].ID = k.ID
	}
	return listContents, nil
}

// StoreContent is used to store a content in the datastore
func (d Data) StoreContent(content structures.ListContent) (int64, error) {
	content.ID = 0
	key := datastore.IncompleteKey(d.listContent, nil)
	key, err := d.client.Put(d.ctx, key, &content)
	if key != nil {
		return key.ID, err
	}
	return 0, err
}

// Delete remove from the datastore the entity of the given key
func (d Data) Delete(key *datastore.Key) error {
	return d.client.DeleteMulti(d.ctx, []*datastore.Key{key})
}

// GetUserKey return the key associated with the user email
func (d Data) GetUserKey(UserMail string) (*datastore.Key, error) {
	q := datastore.NewQuery(d.users).
		Filter("Email =", UserMail).
		Limit(1).
		KeysOnly()
	keys, err := d.client.GetAll(d.ctx, q, nil)
	if err != nil || len(keys) == 0 {
		return nil, errors.New("Could no fetch users")
	}
	return keys[0], nil
}

// GetStoreKey return the key attached to a particuliar store
func (d Data) GetStoreKey(store structures.Store) (*datastore.Key, error) {
	q := datastore.NewQuery(d.stores).
		Filter("Name =", store.Name).
		Filter("Address =", store.Address).
		Filter("Country =", store.Country).
		Filter("Zip =", store.Zip).
		Limit(1).
		KeysOnly()
	keys, err := d.client.GetAll(d.ctx, q, nil)
	if err != nil || len(keys) == 0 {
		return nil, errors.New("Could no fetch stores")
	}
	return keys[0], nil
}

// ScoreStore update the score of a store
func (d Data) ScoreStore(store structures.Store) error {
	key, err := d.GetStoreKey(store)
	if err != nil {
		return err
	}
	var s structures.Store
	if err := d.client.Get(d.ctx, key, &s); err != nil {
		return err
	}
	s.Score = math.ScoreStore(s.Score, store.Score, s.ScoreNb)
	s.ScoreNb++
	_, err = d.client.Put(d.ctx, key, &s)
	return err
}

// UpdateUser update user fields
func (d Data) UpdateUser(user structures.User) error {
	key, err := d.GetUserKey(user.Email)
	if err != nil {
		return err
	}
	var u structures.User
	if err := d.client.Get(d.ctx, key, &u); err != nil {
		return err
	}
	u = structures.MergeUser(user, u)
	_, err = d.client.Put(d.ctx, key, &u)
	return err
}

// DeleteStore remove a store from storage
func (d Data) DeleteStore(ID int64) error {
	return d.Delete(datastore.IDKey(d.stores, ID, nil))
}

// UpdateStore update store fields
func (d Data) UpdateStore(store structures.Store) error {
	key := datastore.IDKey(d.stores, store.ID, nil)
	var s structures.Store
	if err := d.client.Get(d.ctx, key, &s); err != nil {
		return err
	}
	s = structures.MergeStore(store, s)
	_, err := d.client.Put(d.ctx, key, &s)
	return err
}
