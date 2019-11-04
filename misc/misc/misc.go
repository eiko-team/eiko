package misc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

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
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(v)
	if err != nil {
		Logger.Printf("\033[31mError\033[0m: '%s'\n", err.Error())
	}
	return err
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

// SplitString return the string s splited with the separator sep and the size
// of result is at least lenRes.
func SplitString(s, sep string, lenRes int) []string {
	var res []string
	for res = strings.Split(s, sep); len(res) < lenRes; {
		res = append(res, "")
	}
	return res
}
