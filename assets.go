package main

import (
	"io/ioutil"

	"github.com/bitmark-inc/bitmark-sdk-go/asset"
	"github.com/bitmark-inc/bitmark-sdk-go/bitmark"
	"github.com/gin-gonic/gin"
)

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
		panic(err)
	}
	if err = a.issue(); err != nil {
		panic(err)
	}
	if err = db.Create(&a).Error; err != nil {
		panic(err)
	}
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
	a.ID, err = asset.Register(params)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) issue() (err error) {
	// issue an asset to the network
	// use quantity to issue multiple asset
	params := bitmark.NewIssuanceParams(a.ID, a.Amount)
	params.Sign(a.Owner.Account())
	bitmarkIDs, err := bitmark.Issue(params)
	if err != nil {
		panic(err)
	}
	for _, id := range bitmarkIDs {
		issue := Issue{
			ID:    id,
			Asset: *a,
		}
		if err := db.Create(&issue).Error; err != nil {
			panic(err)
		}
	}
	return nil
}
