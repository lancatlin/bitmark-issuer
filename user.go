package main

import (
	"github.com/bitmark-inc/bitmark-sdk-go/account"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func signUp(c *gin.Context) {
	// 註冊帳號
	// 為使用者創建錢包
	// Handle Ajax
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
		if err.Error() == "UNIQUE constraint failed: users.email" {
			c.String(401, "email already used")
		} else {
			panic(err)
		}
	}
	// 註冊成功
	login(c)
}
