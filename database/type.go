package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ExtraArtist struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	OtherName *string `json:"other_name"`
}

type ExtraArtists []ExtraArtist

func (s ExtraArtists) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (v *ExtraArtists) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch value := value.(type) {
	case string:
		return json.Unmarshal([]byte(value), &v)
	case []byte:
		return json.Unmarshal(value, &v)
	default:
		return errors.New(fmt.Sprintf("unsupported type %T", v))
	}
}
