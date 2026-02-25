package value_object

import (
	"github.com/myproject/api/domain/exception"
	"strings"
)

type Gender string

var GenderMale = Gender("male")
var GenderFemale = Gender("female")

func ParseGender(gender string) (Gender, error) {
	validGenders := map[string]Gender{
		GenderMale.String():   GenderMale,
		GenderFemale.String(): GenderFemale,
	}

	g := strings.ToLower(gender)
	if value, exists := validGenders[g]; exists {
		return value, nil
	}

	return "", exception.New("invalid gender; choose a correct value (male or female)")
}

func (g Gender) String() string {
	return string(g)
}
