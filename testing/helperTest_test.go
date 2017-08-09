package testing

// Source: https://speakerdeck.com/mitchellh/advanced-testing-with-go

import (
	"os"
	"testing"
)

func testChdir(t *testing.T, dir string) func() {
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("err: %s", err)
	}

	return func() { os.Chdir(old) }
}

func TestThing(t *testing.T) {
	defer testChdir(t, "/tmp")
}
