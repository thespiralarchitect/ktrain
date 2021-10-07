package tokens

import (
	"math/rand"
	"strconv"
	"time"
)

func CreateToken(userId int64, username string, birthday time.Time, created_at time.Time) string {
	token := birthday.String() + time.Now().String() + created_at.String() + username + strconv.Itoa(int(userId)) + strconv.Itoa(rand.Intn(100))
	return token
}
