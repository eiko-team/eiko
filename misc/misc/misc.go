package misc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Ompluscator/dynamic-struct"
	jwt "github.com/dgrijalva/jwt-go"

	"eiko/misc/hash"
	"eiko/misc/structures"
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

// ParseStruct generic function to parse request body, extract it's content,
// create a struct, fill it and return it
func ParseStruct(r *http.Request) (interface{}, error) {
	if r == nil || r.Body == nil {
		return nil, fmt.Errorf("No Body")
	}

	// mapping request's body to map
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	bodyMap := make(map[string]interface{})
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	// creating new struct
	s := dynamicstruct.NewStruct()
	for key, value := range bodyMap {
		json := strings.ToLower(key)
		firestore := strings.Title(json)
		tag := fmt.Sprintf(`json:"%s" firestore:"%s"`, json, firestore)
		switch t := reflect.ValueOf(value); t.Kind() {
		case reflect.Int:
			v, err := Atoi(value.(string))
			if err != nil {
				return nil, err
			}
			s.AddField(firestore, v, tag)
		case reflect.String:
			s.AddField(firestore, value, tag)
		}
	}
	i := s.Build().New()
	req, err := http.NewRequest("POST", "/",
		strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	err = ParseJSON(req, &i)
	if err != nil {
		return nil, err
	}
	Logger.Printf("%+v", i)
	return i, nil
}

// LogRequest logs a *http.Request using the Logger
func LogRequest(r *http.Request) {
	if r == nil {
		return
	}

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		Logger.Println(err)
	}
	Logger.Println(fmt.Sprintf("%q", requestDump))
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
func TokenToUser(tokenStr string) (structures.User, error) {
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
		return structures.User{
			Email: claims["email"].(string),
		}, nil
	}
	return structures.User{}, err
}

// ValidateToken check if the token is valid
func ValidateToken(token string) bool {
	_, err := TokenToUser(token)
	return err == nil
}

// Atoi is a wrapper around strconv.Atoi
func Atoi(value string) (int, error) {
	i64, err := strconv.ParseInt(value, 10, 0)
	return int(i64), err
}
