package common

// BoolPtr values
var (
	BoolPtrTrue  = BoolPtr(true)
	BoolPtrFalse = BoolPtr(false)
)

// BoolPtr converts the bool value specified as a parameter
// into a bool pointer.
func BoolPtr(b bool) *bool {
	return &b
}
