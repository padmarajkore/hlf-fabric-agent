package utils

// ContainsIgnoreCase checks if substr is in s, case-insensitive
func ContainsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || (len(s) > 0 && (StringIndexFold(s, substr) >= 0)))
}

// StringIndexFold is like strings.Index but case-insensitive
func StringIndexFold(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if EqualFold(s[i:i+len(substr)], substr) {
			return i
		}
	}
	return -1
}

// EqualFold is like strings.EqualFold
func EqualFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c1 := s[i]
		c2 := t[i]
		if c1 == c2 {
			continue
		}
		// To lower-case ASCII
		if 'A' <= c1 && c1 <= 'Z' {
			c1 += 'a' - 'A'
		}
		if 'A' <= c2 && c2 <= 'Z' {
			c2 += 'a' - 'A'
		}
		if c1 != c2 {
			return false
		}
	}
	return true
}
