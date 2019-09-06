package cellouuidgen

import "testing"

func Test_Uuidgen(t *testing.T) {
	_, err := Uuidgen()
	if err != nil {
		t.Fatal("Failed to generate uuid!")
	}
}
