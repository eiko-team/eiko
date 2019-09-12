package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// saltPassword is to generate to byte array to use as the password to hash.
// In there, pass, salt, and some otherthings are shuffled.
func saltPassword(pass, salt string) []byte {
	return []byte(salt + pass + salt + pass)
}

// Hash hashes the password with the salt and return the hash as a string by
// using the bcrypt algorithm.
func Hash(pass, salt string) (string, error) {
	// https://crackstation.net/hashing-security.htm
	h := saltPassword(pass, salt)
	res, err := bcrypt.GenerateFromPassword(h, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// CompareHash compares a password and the salt with the hash.
// return an error if the password/salt is not correct
func CompareHash(h1, pass, salt string) error {
	h := saltPassword(pass, salt)
	return bcrypt.CompareHashAndPassword([]byte(h1), h)
}