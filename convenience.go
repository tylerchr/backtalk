package backtalk

import (
	"strconv"
	"strings"
	"time"
)

// ParseTime is a convenience method for converting a Slack timestamp to a
// native time.Time value.
func ParseTime(t string) time.Time {

	parts := strings.Split(t, ".")
	sec, err1 := strconv.ParseInt(parts[0], 10, 64)
	if err1 != nil {
		panic(err1)
	}
	subsec, err2 := strconv.ParseFloat("0."+parts[1], 64)
	if err2 != nil {
		panic(err1)
	}

	ns := int64(subsec * float64(time.Second))

	return time.Unix(sec, ns)

}
