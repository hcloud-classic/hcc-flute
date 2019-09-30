package checkroot

import "testing"

func Test_Root(t *testing.T) {
	if !CheckRoot() {
		t.Fatal("Failed to get root permission!")
	}
}
