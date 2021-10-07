package errors

import (
	"errors"

	"gorm.io/gorm"
)

func IsDataNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
