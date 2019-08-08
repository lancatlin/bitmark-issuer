package main

import (
	"github.com/pborman/uuid"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// sessions stores "session id: user id"
var sessions = make(map[string]uint)

func getUser(c *gin.Context) User {
	// Get the user from session
	// The session id is stored in cookie
	sessID, err := c.Cookie("session")
	if err != nil {
		return User{IsLogin: false}
	}
	if _, ok := sessions[sessID]; !ok {
		return User{IsLogin: false}
	}
	var user User 
	if err := db.First(&user, sessions[sessID]).Error; err != nil {
		panic(err)
	}
	user.IsLogin = true
	return user
}

func login(c *gin.Context) {
	// 登入帳號
	var user User
	if getUser(c).IsLogin {
		// Already logged in
	}
	if err := c.ShouldBind(&user); err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}
	if err := db.Where(User{Email: user.Email}).First(&user).Error; err != nil {
		// user not found
		c.Status(401)
		return
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(user.PlainPassword)); err != nil {
		// Password not correct
		c.Status(401)
		return
	}
	// logged in
	sessID := uuid.NewRandom().String()
	sessions[sessID] = user.ID
	// response success
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "session",
		Value: sessID,
		Expires: time.Now().Add(time.Minute * 30),
	})
}