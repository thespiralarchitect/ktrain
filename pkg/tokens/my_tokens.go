package tokens

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type MyClaims struct {
	jwt.StandardClaims
	UserID int64
}

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
func GetJWT(userId int64) (string, error) {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		UserID: userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(viper.GetString("key.myKey")))
	if err != nil {
		return "", errors.New("couldn't Signedstring")
	}
	return ss, nil
}
func ParseJWT(cookieValue string) (*jwt.Token, error) {
	afterVertificationToken, err := jwt.ParseWithClaims(cookieValue, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("SOMEONE TRIED TO HACK changed signing method")
		}
		return []byte(viper.GetString("key.myKey")), nil
	})
	if err != nil {
		return nil, errors.New("couldn't parse")
	}
	return afterVertificationToken, nil
}
