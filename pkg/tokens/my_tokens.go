package tokens

import (
	"strconv"
	"time"
)

func CreateToken(userId int64, username string, birthday time.Time, created_at time.Time) string {
	token := birthday.String() + created_at.String() + username + strconv.Itoa(int(userId))
	return token
}
