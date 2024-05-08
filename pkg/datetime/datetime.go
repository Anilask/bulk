package datetime

import (
	"time"

	"bulk/constant"
)

// GetTime will return time as string
func GetTime() string {
	loc, _ := time.LoadLocation(constant.AsiaJakartaTimeZone)
	return time.Now().In(loc).Format(constant.FormatTimeYYYYMMDDHHMMSS)
}
