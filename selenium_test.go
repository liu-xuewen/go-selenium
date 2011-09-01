package selenium

import (
	"testing"
)

func TestStatus(t *testing.T) {
	wd, err := New(nil, "", nil)
	if err != nil {
		t.Error(err)
	}

	status, err := wd.Status()
	if err != nil {
		t.Error(err)
	}
}
