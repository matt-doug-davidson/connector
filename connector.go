package connector

import (
	"strings"
	"time"
)

func after(target string, after string) string {
	// Get substring after a string.
	pos := strings.LastIndex(target, after)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(after)
	if adjustedPos >= len(target) {
		return ""
	}
	return target[adjustedPos:len(target)]
}

// FormatESPTime format date/time to UTC value
func FormatESPTime(timestring string) string {
	// Get the time zone. It is going to be added to the
	// time string before converting.
	now := time.Now()
	timeZone, _ := now.Zone()

	// Make sure the milliseconds have 3 digits
	decimal := after(timestring, ".")
	switch len(decimal) {
	case 0:
		timestring += ".000"
	case 1:
		timestring += "00"
	case 2:
		timestring += "0"
	default:
		return "Error"
	case 3:
	}

	// Remove the T and/or _ between the date and time
	timeStr := strings.Replace(timestring, "T", " ", 1)
	timeStr = strings.Replace(timeStr, "_", " ", 1)
	// Add the time zone.
	datetime := timeStr + " " + timeZone
	// Convert from current location
	t, err := time.Parse("2006-01-02 15:04:05.000 MST", datetime)
	if err != nil {
		return ""
	}

	// Set the location to UTC
	loc, _ := time.LoadLocation("UTC")
	// Determine the current time in the UTC timezone
	utcTime := t.In(loc)
	// Convert to ESP format.
	utcTimeStr := utcTime.Format("2006-01-02T15:04:05.000Z")
	return utcTimeStr
}
