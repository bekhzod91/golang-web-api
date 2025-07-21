package value_object

import (
	"github.com/hzmat24/api/domain/exception"
)

type Rating int64

func ParseRating(rating int64) (Rating, error) {
	if rating < 0 || rating > 5 {
		return 0, exception.New("rating must be between 0 and 5")
	}

	return Rating(rating), nil
}

func (r Rating) Int64() int64 {
	return int64(r)
}
