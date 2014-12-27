package packet

import "testing"

func Test_validateUNSUBACKBytes_ptypeErr(t *testing.T) {
	if err := validateUNSUBACKBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}
