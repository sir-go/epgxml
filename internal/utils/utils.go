package utils

func CutStr(s string, maxLen int) string {
	r := []rune(s)
	if len(r) > maxLen {
		return string(r[:maxLen])
	} else {
		return s

	}
}
