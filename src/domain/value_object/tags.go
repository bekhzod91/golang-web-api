package value_object

import (
	"encoding/json"
	"errors"
)

type Tags []string

func ParseTags(tags []string) (Tags, error) {
	return tags, nil
}

func (p *Tags) Value() []string {
	var items []string
	for _, val := range *p {
		items = append(items, val)
	}

	return items
}

func (p *Tags) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &p)
}
