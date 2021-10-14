package tokens

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const MYKEY = "tovinhtuan"

type myClaims struct {
	jwt.StandardClaims
	Fullname string
	Username string
	Password string
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
func GetJWT(fullname string, username string) (string, error) {
	claims := myClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix(),
		},
		Fullname: fullname,
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(MYKEY))
	if err != nil {
		return "", errors.New("couldn't Signedstring")
	}
	return ss, nil
}
func ParseJWT(cookieValue string) (*jwt.Token, error) {
	afterVertificationToken, err := jwt.ParseWithClaims(cookieValue, &myClaims{}, func(beforeVertificationToken *jwt.Token) (interface{}, error) {
		if beforeVertificationToken.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("SOMEONE TRIED TO HACK changed signing method")
		}
		return []byte(MYKEY), nil
	})
	if err != nil {
		return nil, errors.New("Couln't parseJWT")
	}
	return afterVertificationToken, nil
}
