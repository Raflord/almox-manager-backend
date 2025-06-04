package utils

import "time"

func ParseDateTime(stringDateTime string) (time.Time, error) {
	dateTime, err := time.Parse(time.DateTime, stringDateTime)
	if err != nil {
		return dateTime, err
	}

	return dateTime, err
}
