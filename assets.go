package main

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/bitmark-inc/bitmark-sdk-go/asset"
	"github.com/bitmark-inc/bitmark-sdk-go/bitmark"
	"github.com/gin-gonic/gin"
)

var (
	// ErrAlreadyRegister mean a file was register in the blockchain
	ErrAlreadyRegister = errors.New("this file already has been register")
)

func getAsset(c *gin.Context) {
	// Everyone who has the url can access to this page
	var a Asset
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&a).Error; err != nil {
		// record not found
		c.String(404, err.Error())
		return
	}
	user := getUser(c)
	page := struct {
		User
		Asset
	}{
		User:  user,
		Asset: a,
	}
	c.HTML(200, "get-asset.html", page)
}

func newAsset(c *gin.Context) {
	user := getUser(c)
	if !user.IsLogin {
		c.Redirect(303, "/login")
		return
	}
	var a Asset
	if err := c.ShouldBind(&a); err != nil {
		panic(err)
	}
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
	a.Owner = user
	if err = a.register(data); err != nil {
		if err == ErrAlreadyRegister {
			page := struct {
				message
				User
			}{
				message: message{
					Title:      "建立資產失敗",
					Content:    "此資產已經存在網路中，無法建立",
					Target:     "/assets/new",
					TargetName: "建立資產",
				},
				User: user,
			}
			c.HTML(403, "msg.html", page)
		} else {
			panic(err)
		}
		return
	}
	if err = a.issue(); err != nil {
		panic(err)
	}
	if err = db.Create(&a).Error; err != nil {
		c.String(500, "db create fatal: %s\n%v", err.Error(), a)
		panic(err)
	}
	c.Redirect(303, "/assets/"+a.ID)
}

func (a *Asset) register(data []byte) (err error) {
	// register an asset to the network but not issue yet
	// need a.issue after register
	acct := a.Owner.Account()
	params, err := asset.NewRegistrationParams(
		a.Name,
		map[string]string{"Creater": a.Owner.Name},
	)
	params.SetFingerprint(data)
	params.Sign(acct)
	chErr := make(chan error, 1)
	chID := make(chan string, 1)
	go func(a *Asset) {
		id, err := asset.Register(params)
		chErr <- err
		chID <- id
	}(a)
	select {
	case err = <-chErr:
		a.AssetID = <-chID
		if err != nil {
			if err.Error() == "[1000] message: invalid parameters reason: asset already registered" {
				return ErrAlreadyRegister
			}
			return err
		}
	case <-time.After(5 * time.Second):
		// Timeout
		// The asset already has been register
		return ErrAlreadyRegister
	}
	a.ID = randomID()
	return nil
}

func (a *Asset) issue() (err error) {
	// issue an asset to the network
	// use quantity to issue multiple asset
	params := bitmark.NewIssuanceParams(a.AssetID, a.Amount)
	params.Sign(a.Owner.Account())
	bitmarkIDs, err := bitmark.Issue(params)
	if err != nil {
		return errors.New("bitmark issue fatal: " + err.Error())
	}
	a.Issues = make([]Issue, a.Amount)
	for i, id := range bitmarkIDs {
		issue := Issue{
			ID:      id,
			AssetID: a.ID,
		}
		a.Issues[i] = issue
	}
	return nil
}
