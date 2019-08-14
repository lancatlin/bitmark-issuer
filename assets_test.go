package main

import (
	"crypto/rand"
	"testing"
)

func TestRegister(t *testing.T) {
	var user User
	if err := db.First(&user).Error; err != nil {
		t.Error("No user is found, exit...")
		return
	}
	var data = make([]byte, 128)
	_, err := rand.Read(data)
	if err != nil {
		t.Error(err)
	}
	a := Asset{
		User:   user,
		Name:   "test asset",
		Amount: 2,
	}
	if err = a.register(data); err != nil {
		t.Error(err)
		return
	}
	if err = a.issue(); err != nil {
		t.Error(err)
		return
	}
	if err := a.register(data); err == nil {
		t.Error("register didn't return error: ", err)
	}
}
