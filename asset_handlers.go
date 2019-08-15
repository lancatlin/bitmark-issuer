package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func assetInfo(c *gin.Context) {
	user := getUser(c)
	var a Asset
	if err := db.Where("id = ?", c.Param("id")).Preload("URL").First(&a).Error; err != nil {
		// not found
		c.String(404, err.Error())
		return
	}
	log.Println(a)
	if user.ID != a.UserID {
		page := struct {
			User
			message
		}{
			User: user,
			message: message{
				Title:   "你沒有權限",
				Content: "你沒有權限瀏覽此頁面，請檢查你的帳號是否正確",
			},
		}
		c.HTML(401, "msg.html", page)
		return
	}
	if _, err := os.Stat("/static/qrcodes/" + a.URL.ID + ".png"); os.IsNotExist(err) {
		if err := genQRCode(a.URL.String()); err != nil {
			panic(err)
		}
	}
	log.Println("fine after genQRCode")
	page := struct {
		User
		Asset
		Remainder int
		Host      string
	}{
		User:  user,
		Asset: a,
		Host:  env.Host,
	}
	c.HTML(200, "asset-info.html", page)
}

func getAsset(c *gin.Context) {
	// Everyone who has the url can access to this page
	user := getUser(c)
	var url URL
	id := c.Param("id")
	if err := db.Where("id = ?", id).Preload("Asset").Preload("Asset.URL").First(&url).Error; err != nil {
		// record not found
		c.String(404, err.Error())
		return
	}
	if time.Now().After(url.ExpireAt) {
		// 連結失效
		page := struct {
			User
			message
		}{
			User: user,
			message: message{
				Title:   "連結失效",
				Content: "此連結已經過期",
			},
		}
		c.HTML(403, "msg.html", page)
		return
	}
	page := struct {
		User
		Asset
	}{
		User:  user,
		Asset: url.Asset,
	}
	c.HTML(200, "get-asset.html", page)
}

func newAsset(c *gin.Context) {
	user := getUser(c)
	var a Asset
	log.Println("fine here")
	a.Name = c.PostForm("asset_name")
	var err error
	a.Amount, err = strconv.Atoi(c.PostForm("amount"))
	if err != nil {
		a.Amount = 1
	}
	log.Println("fine after binding")
	header, err := c.FormFile("file")
	if err != nil {
		panic(err)
	}
	file, err := header.Open()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	log.Println("fine after loading file")
	a.ID = randomID()
	a.User = user
	if err = a.register(data); err != nil {
		if err == ErrAlreadyRegister {
			page := struct {
				message
				User
			}{
				message: message{
					Title:      "建立資產失敗",
					Content:    "此資產已經存在網路中，無法建立",
					Target:     "/new",
					TargetName: "建立資產",
				},
				User: user,
			}
			c.HTML(403, "msg.html", page)
		} else if err == ErrTimeout {
			c.String(500, err.Error())
			return
		} else {
			c.String(500, err.Error())
			return
		}
		return
	}
	log.Println("fine after register")
	if err = a.issue(); err != nil {
		panic(err)
	}
	log.Println("fine after issue")
	expires, err := strconv.Atoi(c.PostForm("expires"))
	if err != nil {
		expires = 1
	}
	a.URL = &URL{
		ID:       randomID(),
		ExpireAt: time.Now().AddDate(0, 0, expires),
	}
	if err = db.Create(&a).Error; err != nil {
		c.String(500, "db create fatal: %s\n%v", err.Error(), a)
		panic(err)
	}
	log.Println("fine after create")
	c.Redirect(303, "/assets/"+a.ID)
}
