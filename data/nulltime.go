package data

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time = v
	case []byte:
		parsed, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		nt.Time = parsed
	case string:
		parsed, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		nt.Time = parsed
	default:
		return fmt.Errorf("cannot scan type %T into NullTime", value)
	}

	nt.Valid = true
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time.Format("2006-01-02 15:04:05"), nil
}
