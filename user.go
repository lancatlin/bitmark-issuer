package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/bitmark-inc/bitmark-sdk-go/account"
)

func signUp(c *gin.Context) {
	// 註冊帳號
	// 為使用者創建錢包
	var user User
	var err error
	if err = c.ShouldBind(&user); err != nil {
		panic(err)
	}
	user.Password, err = bcrypt.GenerateFromPassword([]byte(user.PlainPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	acct, err := account.New()
	user.Wallet = acct.Seed()
	if err != nil {
		panic(err)
	}
	if err := db.Create(&user).Error; err != nil {
		// something went wrong
		// 可能是帳號重複
		panic(err)
	}
	// 註冊成功
	page := Message{"註冊成功", "註冊成功，登入後即可使用", "/login", "登入"}
	c.HTML(200, "msg.html", page)
}