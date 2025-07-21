package helper

import (
	"errors"
	"strconv"
)

func Atoi(s string) (int64, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return int64(id), nil
}
