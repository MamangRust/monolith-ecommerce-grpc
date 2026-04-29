package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ParseDate parses a string in "2006-01-02" format to pgtype.Date
func ParseDate(dateStr string) (pgtype.Date, error) {
	if dateStr == "" {
		return pgtype.Date{Valid: false}, nil
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return pgtype.Date{Valid: false}, fmt.Errorf("invalid date format: %v", err)
	}

	return pgtype.Date{Time: t, Valid: true}, nil
}

// ParseTime parses a string in "15:04:05" format to pgtype.Time
func ParseTime(timeStr string) (pgtype.Time, error) {
	if timeStr == "" {
		return pgtype.Time{Valid: false}, nil
	}

	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return pgtype.Time{Valid: false}, fmt.Errorf("invalid time format: %v", err)
	}

	// pgtype.Time represents time in microseconds since midnight
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	microseconds := t.Sub(midnight).Microseconds()

	return pgtype.Time{Microseconds: microseconds, Valid: true}, nil
}
