package main

import (
	"crypto/rand"
	"encoding/base64"
)

func randomID() string {
	// Generate a random string
	// code by base64
	// length 8
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
