package entity

import (
	"database/sql/driver"
	"encoding/json"
)

type DecimalArray []float64

func (a *DecimalArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *DecimalArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil // or return an error, depending on your requirements
	}

	var array []float64
	if err := json.Unmarshal(bytes, &array); err != nil {
		return err
	}

	*a = array
	return nil
}
