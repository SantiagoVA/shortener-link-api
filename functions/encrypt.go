package functions

import "github.com/alexedwards/argon2id"

func Encrypt(k string) (string, error) {
	hash, err := argon2id.CreateHash(k, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}
