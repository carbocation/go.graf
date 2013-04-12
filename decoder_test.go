package main

import (
	"testing"
)

//This tests the custom convertPassword decoder we've created that uses
// Gorilla Web Toolkit's Schema package.
func TestConvertPassword(t *testing.T) {
	expected := Password("PASS")

	user := new(User)

	input := map[string][]string{
		"Password": {"PASS"},
	}

	decoder.Decode(user, input)

	if user.Password != expected {
		t.Errorf("Password was %s, expected %s", user.Password, expected)
	}
}
