package utils

import (
	"fmt"
)

// HumanTimeSecond 格式化秒数
func HumanTimeSecond(s int64, separator ...string) string {
	const (
		Minute = 60
		Hour = 60 * Minute
		Day = 24 * Hour
	)
	var d, h, m int64

	d = s / Day
	s = s % Day
	h = s / Hour
	s = s % Hour
	m = s / Minute
	s = s % Minute

	var sep string
	if len(separator) == 1 {
		sep = separator[0]
	}

	if d > 0 {
		return fmt.Sprintf("%dd%s%dh%s%dm%s%ds", d, sep, h, sep, m, sep, s)
	} else if h > 0 {
		return fmt.Sprintf("%dh%s%dm%s%ds", h, sep, m, sep, s)
	} else if m > 0 {
		return fmt.Sprintf("%dm%s%ds", m, sep, s)
	}
	return fmt.Sprintf("%ds", s)
}

// HumanByteCount 美化字节显示, 默认为国际单位制制定的十进制标准(SI), isIEC 为 true 时为国际电工委员会制定的二进制标准(IEC)
func HumanByteCount(b int64, isIEC ...bool) string {
	var unit int64
	if len(isIEC) == 1 && isIEC[0] {
		unit = 1024
	} else {
		unit = 1000
	}

	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	if unit == 1024 {
		return fmt.Sprintf("%.1f %ciB",
			float64(b)/float64(div), "KMGTPE"[exp])
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}