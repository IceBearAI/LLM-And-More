package util

import (
	"fmt"
	"strings"
	"time"
)

// DurationPrecision 定义了持续时间的精度选项
type DurationPrecision int

const (
	PrecisionDays DurationPrecision = iota
	PrecisionHours
	PrecisionMinutes
	PrecisionSeconds
	PrecisionMilliseconds
)

func FormatDuration(seconds float64, precision DurationPrecision) string {
	duration := time.Duration(seconds * float64(time.Second))
	var parts []string
	if duration <= 0 {
		return ""
	}
	days := duration / (24 * time.Hour)
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if precision == PrecisionDays {
		return strings.Join(parts, "")
	}

	duration -= days * 24 * time.Hour
	hours := duration / time.Hour
	if hours > 0 || len(parts) > 0 { // 即使小时为0，如果已有天数也显示
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if precision == PrecisionHours {
		return strings.Join(parts, "")
	}

	duration -= hours * time.Hour
	minutes := duration / time.Minute
	if minutes > 0 || len(parts) > 0 { // 即使分钟为0，如果已有小时也显示
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if precision == PrecisionMinutes {
		return strings.Join(parts, "")
	}

	duration -= minutes * time.Minute
	secs := duration / time.Second
	if secs > 0 || len(parts) > 0 { // 即使秒为0，如果已有分钟也显示
		parts = append(parts, fmt.Sprintf("%ds", secs))
	}
	if precision == PrecisionSeconds {
		return strings.Join(parts, "")
	}

	duration -= secs * time.Second
	millis := duration / time.Millisecond
	if millis > 0 || len(parts) > 0 { // 即使毫秒为0，如果已有秒也显示
		parts = append(parts, fmt.Sprintf("%dms", millis))
	}
	return strings.Join(parts, "")
}
