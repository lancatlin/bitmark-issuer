package main

import "testing"

func TestRandomString(t *testing.T) {
	str := randomID()
	if len([]byte(str)) != 8 {
		t.Error("Length of randomID is not 8")
	}
	t.Log(str)
}
