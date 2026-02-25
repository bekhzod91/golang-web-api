package value_object

import (
	"github.com/myproject/api/domain/exception"
	"strings"
)

type Status string

const StatusActive = Status("active")
const StatusInactive = Status("inactive")

func ParseStatus(status string) (Status, error) {
	s := strings.ToLower(status)
	if s == StatusActive.String() || s == StatusInactive.String() {
		return Status(s), nil
	}

	return "", exception.New("invalid status choose correct value (active, inactive)")
}

func (s Status) String() string {
	return string(s)
}
