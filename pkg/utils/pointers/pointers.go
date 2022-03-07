package pointers

// Returns a pointer to a string - nil if string is empty
func String(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func Bool(b bool) *bool {
	return &b
}
