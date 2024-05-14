package web

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncrypt(T *testing.T) {
	password := "hello@world"
	encrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		T.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypt, []byte(password))
	assert.NoError(T, err)
}
