package common

// BoolPtr values
var (
	BoolPtrTrue  = BoolPtr(true)
	BoolPtrFalse = BoolPtr(false)
)

// BoolPtr converts a bool value into a bool pointer.
func BoolPtr(b bool) *bool {
	return &b
}
