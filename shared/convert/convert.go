package convert

// NullableString converts an empty string to a nil pointer.
func NullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// StringPtr returns a pointer to s.
func StringPtr(s string) *string {
	return &s
}
