package misc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/hash"
	"github.com/eiko-team/eiko/misc/log"
	"github.com/eiko-team/eiko/misc/structures"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// TokensKey is the secret key for all tokens.
	// warning: On server reboot, all token are invalidated
	TokensKey = []byte(hash.GenerateKey(666))

	// ExpToken number of days before the token expire
	ExpToken = time.Duration(14)

	// Logger used to log output
	Logger = log.New(os.Stdout, "Misc: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// ParseJSON generic function to parse request body, extract it's content and
// fill the struct
func ParseJSON(r *http.Request, v interface{}) error {
	if r == nil || r.Body == nil {
		return fmt.Errorf("No Body")
	}

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// DumpRequest return the request in a string format
func DumpRequest(r *http.Request) string {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		Logger.Println(err)
	}
	return string(requestDump)
}

// LogRequest logs a *http.Request using the Logger
func LogRequest(r *http.Request) {
	if r == nil {
		return
	}
	Logger.Println(DumpRequest(r))
}

// UserToToken convert the user information to a valid token
func UserToToken(u structures.User) (string, error) {
	return jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), jwt.MapClaims{
		"email": u.Email,
		"exp":   time.Now().Add(time.Hour * 24 * ExpToken).Unix(),
		"iat":   time.Now().Unix(),
		// "nbf":   time.Now().Add(time.Second * 1).Unix(),
	}).SignedString(TokensKey)
}

// TokenToUser convert the Token to user's information
func TokenToUser(d data.Data, tokenStr string) (structures.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Bad method: %v", token.Header["alg"])
		}
		return TokensKey, nil
	})
	if err != nil {
		return structures.User{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return d.GetUser(claims["email"].(string))
	}
	return structures.User{}, err
}

// ValidateToken check if the token is valid
func ValidateToken(d data.Data, token string) bool {
	_, err := TokenToUser(d, token)
	return err == nil
}

// Atoi is a wrapper around strconv.Atoi
func Atoi(value string) (int, error) {
	i64, err := strconv.ParseInt(value, 10, 0)
	return int(i64), err
}

// SplitString return the string s splited with the separator sep and the size
// of result is at least lenRes.
func SplitString(s, sep string, lenRes int) []string {
	var res []string
	for res = strings.Split(s, sep); len(res) < lenRes; {
		res = append(res, "")
	}
	return res
}

// NormalizeName is used to normalize a name
func NormalizeName(name string) string {
	return strings.ToLower(strings.Replace(name, "%20", " ", -1))
}

// NormalizeConsumable is used to normalize a consumable
func NormalizeConsumable(c structures.Consumable) structures.Consumable {
	c.Name = NormalizeName(c.Name)
	return c
}

// NormalizeQuery is used to normalize a query
func NormalizeQuery(q structures.Query) structures.Query {
	q.Query = NormalizeName(q.Query)
	if q.Limit < 1 || q.Limit > 20 {
		q.Limit = 20
	}
	return q
}

// MergeUser Merge two user into one:
// i2 fill empty fields of i1
func MergeUser(i1, i2 structures.User) structures.User {
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
	if !i1.SBio {
		i1.SBio = i2.SBio
	}
	if !i1.SVegan {
		i1.SVegan = i2.SVegan
	}
	if !i1.SHalal {
		i1.SHalal = i2.SHalal
	}
	if !i1.SCasher {
		i1.SCasher = i2.SCasher
	}
	if !i1.SSodium {
		i1.SSodium = i2.SSodium
	}
	if !i1.SEgg {
		i1.SEgg = i2.SEgg
	}
	if !i1.SPenut {
		i1.SPenut = i2.SPenut
	}
	if !i1.SCrustace {
		i1.SCrustace = i2.SCrustace
	}
	if !i1.SGluten {
		i1.SGluten = i2.SGluten
	}
	if !i1.SDiabetic {
		i1.SDiabetic = i2.SDiabetic
	}
	if len(i1.SBoycott) != len(i2.SBoycott) {
		// TODO: Implement an array merging system
	}
	return i1
}

// MergeStore Merge two user into one:
// i2 fill empty fields of i1
func MergeStore(i1, i2 structures.Store) structures.Store {
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

// MergeConsumable Merge two consumable into one:
// C2 fill empty fields of C1
func MergeConsumable(c1, c2 structures.Consumable) structures.Consumable {
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
