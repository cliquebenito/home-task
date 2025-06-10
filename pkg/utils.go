package pkg

import (
	"fmt"
	"time"
)

// parseFlexibleTime attempts to parse a time string in multiple formats
func parseFlexibleTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.000000",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", value)
}
