package api

import (
	"strings"
	"testing"
)

func TestReadBytesLimit(t *testing.T) {
	got, err := readBytesLimit(strings.NewReader("abcde"), 5)
	if err != nil || string(got) != "abcde" {
		t.Fatalf("at limit: got %q err=%v", got, err)
	}
	if _, err := readBytesLimit(strings.NewReader("abcdef"), 5); err == nil {
		t.Fatal("above limit must be rejected")
	}
}
