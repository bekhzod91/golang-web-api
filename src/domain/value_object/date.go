package value_object

import (
	"fmt"
	"github.com/myproject/api/domain/exception"
	"time"
)

type Date time.Time

// ParseDate validates and returns a new Date object
func ParseDate(value string) (Date, error) {
	// Parse the string into a time.Time object
	parsedTime, err := time.Parse(time.DateOnly, value)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return Date{}, exception.New("invalid date format YYYY-mm-dd")
	}
	return Date(parsedTime), nil
}

func (d Date) String() string {
	return time.Time(d).Format(time.DateOnly) // Format it as needed
}

func (d Date) Time() time.Time {
	return time.Time(d)
}

//func (d Date) UnmarshalJSON(data []byte) error {
//	var result string
//
//	err := json.Unmarshal(data, &result)
//	if err != nil {
//		return fmt.Errorf("date error unmarshaling JSON: %w", err)
//	}
//	parsedDate, err := ParseDate(result)
//	if err != nil {
//		return err
//	}
//	*d = parsedDate
//	return nil
//}
