package hb_test

import (
	"testing"

	hb "github.com/parf/homebase-go-lib"
)

func TestVersion(t *testing.T) {
	if hb.Version == "" {
		t.Error("Version should not be empty")
	}
}
