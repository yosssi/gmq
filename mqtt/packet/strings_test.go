package packet

import "testing"

func Test_appendLenStr(t *testing.T) {
	got := appendLenStr([]byte{}, []byte{0x01})

	want := []byte{0x00, 0x01, 0x01}

	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Errorf("got => %v, want => %v", got, want)
	}
}
