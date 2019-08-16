package main

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

func init() {
	// init dir
	if _, err := os.Stat("static/qrcodes"); os.IsNotExist(err) {
		if err := os.Mkdir("static/qrcodes", 0755); err != nil {
			panic(err)
		}
	}
}

func randomID() string {
	// Generate a random string
	// code by base64
	// length 8
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func genQRCode(filename, url string) (err error) {
	return qrcode.WriteFile(url, qrcode.Medium, 256, filename)
}
