package research

import (
	"os"

	"github.com/eiko-team/eiko/misc/log"
	"github.com/eiko-team/eiko/misc/research/helper"
	"github.com/eiko-team/eiko/misc/structures"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "research: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

// Research container for all search relative variables
// initiated inside InitReseach(...)
type Research struct {
	// client is used to take advantage of the datastore api
	client *search.Client
	// consumables log name inside the datastore
	consumables *search.Index
}

// InitReseach return an initialised Research struct
func InitReseach() Research {
	var r Research

	appIDStr := "SEARCH_APP_ID"
	appID := os.Getenv(appIDStr)
	apiKeyStr := "SEARCH_API_KEY"
	apiKey := os.Getenv(apiKeyStr)
	if appID == "" || apiKey == "" {
		Logger.Fatalf("please set: '%s' and '%s'", appIDStr, apiKeyStr)
	}
	r.client = search.NewClient(appID, apiKey)

	r.consumables = r.client.InitIndex("Consumables")
	return r
}

// StoreConsumable stores a consumable in the seach database for later search
func (r Research) StoreConsumable(c structures.Consumable) error {
	_, err := r.consumables.SaveObjects(reseachhelper.ConsumableToSearchable(c),
		opt.AutoGenerateObjectIDIfNotExist(true))
	return err
}

// SearchConsumable search a consumable in the seach database
func (r Research) SearchConsumable(query string) ([]int64, error) {
	q, err := r.consumables.Search(query)
	if err != nil {
		return nil, err
	}
	res := make([]int64, len(q.Hits))
	for i, val := range q.Hits {
		res[i] = int64(val["ID"].(float64))
	}
	return res, nil
}
