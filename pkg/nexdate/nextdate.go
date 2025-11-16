package nexdate

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("параметр repeat не может быть пустым")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("неверный формат даты dstart")
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", errors.New("неверный формат правила")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("не указан интервал в днях")
		}
		interval, err := parseInterval(parts[1])
		if err != nil {
			return "", err
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	default:
		return "", errors.New("неподдерживаемый формат")
	}
}

func parseInterval(intervalStr string) (int, error) {
	var interval int
	_, err := fmt.Sscan(intervalStr, &interval)
	if err != nil || interval <= 0 || interval > 400 {
		return 0, errors.New("неверный интервал")
	}
	return interval, nil
}
