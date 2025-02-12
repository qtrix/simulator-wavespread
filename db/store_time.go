package db

import (
	"database/sql/driver"
	"time"

	"github.com/pkg/errors"
)

type StoreTime time.Time

func (t *StoreTime) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("time in invalid format, use RFC3339")
	}
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*t = StoreTime(parsed)
	return nil
}

func (t StoreTime) Value() (driver.Value, error) {
	return time.Time(t).Format(time.RFC3339), nil
}
