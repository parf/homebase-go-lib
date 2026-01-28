package homebase_test

import (
	"testing"

	"github.com/parf/homebase-go-lib"
)

func TestVersion(t *testing.T) {
	if homebase.Version == "" {
		t.Error("Version should not be empty")
	}
}
