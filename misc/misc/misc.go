package misc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"cloud.google.com/go/datastore"
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

// UniqEmail is used to find if a email is already used in the datastore
func UniqEmail(ctx context.Context, Users, UserMail string,
	client *datastore.Client) error {
	// Finding if the email is unique
	var users []structures.User
	q := datastore.NewQuery(Users).Filter("Email =", UserMail).Limit(1)
	if _, err := client.GetAll(ctx, q, &users); err != nil {
		return err
	}
	if len(users) != 0 {
		return errors.New("user found")
	}
	return nil
}

// LogRequest logs a *http.Request using the Logger
func LogRequest(r *http.Request) {
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
