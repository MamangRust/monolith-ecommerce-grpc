package repository

func IntPtrToInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func stringPtr(v string) *string {
	return &v
}

func int32Ptr(v int) *int32 {
	res := int32(v)
	return &res
}
