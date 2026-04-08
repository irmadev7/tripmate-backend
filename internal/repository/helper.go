package repository

import (
	"errors"

	"github.com/lib/pq"
)

const uniqueViolationCode = "23505"

func IsDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == uniqueViolationCode
	}
	return false
}
