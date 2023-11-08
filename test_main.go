package main

import (
	"testing"
)

func TestGenerateHCL(t *testing.T) {
	input := "hi"
	want := "hi"
	if !want.MatchString(input) {
		t.Errorf("want %s, got %s", want, input)
}
