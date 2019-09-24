package hash

import (
	"math/rand"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var (
	keyPossibilities = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789&(-_)=+$*!:;,?./<>[]{}%@`^\\~")

	// salt salt for the passwrd hashing
	// salt = GenerateKey(666)
	salt = os.Getenv("SALT")
)

// GenerateKey generate a n th long key
func GenerateKey(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = keyPossibilities[rand.Intn(len(keyPossibilities))]
	}
	return string(b)
}

// saltPassword is to generate to byte array to use as the password to hash.
// In there, pass, salt, and some otherthings are shuffled.
func saltPassword(pass string) []byte {
	return []byte(pass + salt + pass + salt + pass)
}

// Hash hashes the password with the salt and return the hash as a string by
// using the bcrypt algorithm.
func Hash(pass string) (string, error) {
	// https://crackstation.net/hashing-security.htm
	h := saltPassword(pass)
	res, err := bcrypt.GenerateFromPassword(h, bcrypt.DefaultCost)
	return string(res), err
}

// CompareHash compares a password and the salt with the hash.
// return if the password is a match for the hash
func CompareHash(hash, pass string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hash),
		saltPassword(pass))
}
