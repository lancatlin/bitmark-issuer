package main

import (
	"io/ioutil"
	"testing"
)

func TestRegister(t *testing.T) {
	var user User
	if err := db.First(&user).Error; err != nil {
		t.Error(err)
	}
	data, err := ioutil.ReadFile("go.sum")
	if err != nil {
		t.Error(err)
	}
	a := Asset{
		Owner:  user,
		Name:   "test asset",
		Amount: 2,
	}
	if err = a.register(data); err != nil {
		t.Error(err)
	}
}
