package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func transfer(c *gin.Context) {
	var url URL
	if err := db.Where("id = ?", c.Param("id")).First(&url).Error; err != nil {
		// not found
		c.String(404, err.Error())
		return
	}
	var issue Issue
	if err := db.Where("asset_id = ?", url.AssetID).Where("receiver = ''").Preload("Asset").Take(&issue).Error; err != nil {
		log.Println(url.AssetID)
		c.String(404, err.Error())
		return
	}
	receiver := c.PostForm("receiver")
	err := issue.transfer(receiver)
	if err == ErrAlreadyTransfer {
		c.String(403, "一個帳戶只能取得一份資產")
		return
	}
	if err != nil {
		c.String(500, err.Error())
		return
	}
	page := struct {
		User
		Asset
	}{
		User:  getUser(c),
		Asset: issue.Asset,
	}
	c.HTML(200, "transfer-success.html", page)
}
