package tokens

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateToken(userId int64, username string, birthday time.Time, created_at time.Time) string {
	token := birthday.String() + time.Now().String() + created_at.String() + username + strconv.Itoa(int(userId)) + strconv.Itoa(rand.Intn(100))
	hasher := md5.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

func HashPassword(password string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Error while generating bcrypt hash from password")
	}
	return bs, nil
}
func ComparePassword(password string, hashedPass []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	if err != nil {
		return errors.New("Invalid password")
	}
	return nil
}
