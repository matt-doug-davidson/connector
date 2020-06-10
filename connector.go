package connector

import (
	"fmt"
	"strings"
	"time"
)

// EspMessage contains the topic and payload
type ConnectorMessage struct {
	Entity   string
	Snapshot MeasurementSnapshot
}

// EspPayload is part of the payload data.
type MeasurementSnapshot struct {
	Datetime     string        `json:"datetime"`
	Measurements []Measurement `json:"values"`
}

// EspValues contains the data for each measurement
type Measurement struct {
	Field      string  `json:"field"`
	Amount     float64 `json:"amount"`
	Attributes string  `json:"attributes"`
}

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
	fmt.Println("decimal: ", decimal)
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
	timeStr = strings.Replace(timestring, "_", " ", 1)
	fmt.Println(timeStr)
	// Add the time zone.
	datetime := timeStr + " " + timeZone
	// Convert from current location
	t, err := time.Parse("2006-01-02 15:04:05.000 MST", datetime)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Set the location to UTC
	loc, _ := time.LoadLocation("UTC")
	// Determine the current time in the UTC timezone
	utcTime := t.In(loc)
	// Convert to ESP format.
	utcTimeStr := utcTime.Format("2006-01-02T15:04:05.000Z")
	return utcTimeStr
}

// Decode decodes the message
func Decode(connectorMessage map[string]interface{}) (string, map[string]interface{}) {

	snapshotMap := make(map[string]interface{})

	msg := connectorMessage["msg"]
	msg1 := msg.(ConnectorMessage)
	snapshot := msg1.Snapshot
	snapshotMap["datetime"] = snapshot.Datetime

	var valuesSlice []interface{}
	for _, value := range snapshot.Measurements {
		field := value.Field
		amount := value.Amount
		attr := value.Attributes
		fmt.Println("\nattr:\n", attr)
		fmt.Printf("%T\n", attr)
		if value.Attributes == "" {
			value := map[string]interface{}{
				"field":  field,
				"amount": amount,
			}
			valuesSlice = append(valuesSlice, value)
		} else {
			value := map[string]interface{}{
				"field":     field,
				"amount":    amount,
				"attribute": attr}
			valuesSlice = append(valuesSlice, value)
		}
	}
	snapshotMap["values"] = valuesSlice
	return msg1.Entity, snapshotMap
}
