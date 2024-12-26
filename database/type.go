package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type FeaturingArtist struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	OtherName *string `json:"other_name"`
}

type FeaturingArtists []FeaturingArtist

func (s FeaturingArtists) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (v *FeaturingArtists) Scan(value any) error {
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
