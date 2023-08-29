package custom

import (
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Datetime represents a datetime.
type Datetime time.Time

// MarshalJSON implements the json.Marshaler interface.
func (d *Datetime) MarshalJSON() ([]byte, error) {
	if d == nil || time.Time(*d).IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`%q`, time.Time(*d).UTC().Format(time.RFC3339))), nil
}

func (d *Datetime) MarshalBSON() ([]byte, error) {
	if d == nil || time.Time(*d).IsZero() {
		return nil, nil
	}
	return []byte(time.Time(*d).UTC().Format(time.RFC3339)), nil
}

func (d *Datetime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if d == nil || time.Time(*d).IsZero() {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(time.Time(*d).UTC().Format(time.RFC3339))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Datetime) UnmarshalJSON(text []byte) error {
	// Remove " from text if present with regex (e.g. "2020-01-01T00:00:00Z" -> 2020-01-01T00:00:00Z)
	reg := regexp.MustCompile(`"(.*)"`)
	text = reg.ReplaceAll(text, []byte("$1"))

	t, err := time.Parse(time.RFC3339, string(text))
	if err != nil {
		return err
	}
	*d = Datetime(t)
	return nil
}

func (d *Datetime) UnmarshalBSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	// Remove all non-alphanumeric characters from the string.
	got := regexp.MustCompile(`[^a-zA-Z0-9-:]`).ReplaceAllString(string(bytes), "")

	// Parse the string into a time.Time.
	t, err := time.Parse(time.RFC3339, got)
	if err != nil {
		return fmt.Errorf("invalid datetime: %s", got)
	}

	// Set the value of the Datetime.
	*d = Datetime(t)
	return nil
}

// Scan implements the sql.Scanner interface.
func (d *Datetime) Scan(src any) error {
	t, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("invalid scan, type %T not supported for %T", src, d)
	}
	*d = Datetime(t)
	return nil
}

// String implements the fmt.Stringer interface.
func (d Datetime) String() string {
	return time.Time(d).Format(time.RFC3339)
}
