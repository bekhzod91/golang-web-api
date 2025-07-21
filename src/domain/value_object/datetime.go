package value_object

import (
	"encoding/json"
	"fmt"
	"github.com/hzmat24/api/domain/exception"
	"strings"
	"time"
)

var (
	ErrInvalidDateTime = exception.New("invalid datetime")
	DateTimeLayout     = time.DateTime
)

type DateTime time.Time

// ParseDateTime validates and returns a new DateTime object
func ParseDateTime(value string) (DateTime, error) {
	// Parse the string into a time.Time object
	value = strings.ReplaceAll(value, "T", " ")

	parsedTime, err := time.Parse(DateTimeLayout, value)
	if err != nil {
		return DateTime{}, ErrInvalidDateTime
	}
	return DateTime(parsedTime), nil
}

func (d *DateTime) String() string {
	return time.Time(*d).Format(DateTimeLayout) // Format it as needed
}

func (d *DateTime) Time() time.Time {
	return time.Time(*d)
}

func (d *DateTime) Date() (Date, error) {
	date, err := ParseDate(time.Time(*d).Format(time.DateOnly))
	if err != nil {
		return date, err
	}
	return date, nil
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	var result string

	err := json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("datetime error unmarshaling JSON: %w", err)
	}

	parsedDatetime, err := ParseDateTime(result)
	if err != nil {
		return err
	}

	*d = parsedDatetime
	return nil
}

func (d *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}
