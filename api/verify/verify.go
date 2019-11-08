package verify

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/misc"
	"github.com/eiko-team/eiko/misc/structures"
)

var (
	// Logger used to log output
	Logger = log.New(os.Stdout, "verify: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

// Email Checks if the email is available
func Email(d data.Data, r *http.Request) (string, error) {
	var i structures.Email
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("2.0.0")
	}

	_, err = d.GetUser(i.UserMail)
	if err == nil {
		return "", errors.New("2.0.1")
	}

	return "{\"available\":\"true\"}", nil
}

// Password Checks the strength of the password
// where 0 <= s <= 4
func Password(d data.Data, r *http.Request) (string, error) {
	var i structures.Password
	err := misc.ParseJSON(r, &i)
	if err != nil {
		return "", errors.New("2.1.0")
	}

	if len(i.Password) == 0 {
		return "{\"strength\":0}", nil
	}

	res := 0
	if len(i.Password) > 9 {
		res++
	}
	var patterns = []struct {
		pattern string
	}{
		{`[a-z]`},                     // abc..
		{`([A-Z]|[0-9])`},             // ABC..., 0123...
		{`([\x21-\x2F]|[\x3A-\x40])`}, // !"#$%&'()*+,-./ || :;<=>?@
		{`([\x5B-\x60]|[\x7B-\x7E])`}, // [\]^_` || {|}~
	}
	for _, tt := range patterns {
		if regexp.MustCompile(tt.pattern).MatchString(i.Password) {
			res++
		}
	}
	return fmt.Sprintf("{\"strength\":%d}", res-1), nil
}
