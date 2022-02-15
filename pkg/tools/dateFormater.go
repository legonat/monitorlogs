package tools

import (
	"fmt"
	"time"
)

const TimeStamp = "2006-01-02 15:04:05 MST"

func FormatDateToUTC(date string) (string, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(t.Format(TimeStamp)), nil
}

func FormatUnixToUTC(date int64) (string) {
	t := time.Unix(date, 0)

	return fmt.Sprint(t.Format(TimeStamp))
}

func FormatDateToUnix(date string) (int64, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

func TimeFix(hour int) time.Duration {
	dif := time.Hour*24 - time.Duration(hour) * time.Hour

	if dif.Hours() > 12 {
		dif = time.Duration(-hour) * time.Hour
	}

	return dif
}