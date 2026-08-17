package util

import "strconv"

// UintString uint 转字符串。
func UintString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
