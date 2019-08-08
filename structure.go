package main

import (
	"io"
	"time"
)

type Asset struct {
	// bitmark asset id
	ID       string `gorm:"primary_key"`
	Name     string `form:"name"`
	Amount   int    `form:"amount"`
	Issues   []Issue
	Owner    User
	UserID   uint
}

type Issue struct {
	ID      string
	Asset   Asset
	AssetID string
	// 如果尚未轉移留空
	Reciever string
	// 如果尚未轉移留空
	TransferredAt time.Time
}

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string `gorm:"unique" form:"email"`
	// 簽署用名稱
	Name          string `form:"name"`
	PlainPassword string `form:"password" gorm:"-"`
	Password      []byte
	// 儲存 account seed
	Wallet  string
	Assets  []Asset
	IsLogin bool `gorm:"-"`
}

func (u User) Account() account.Account {
	acct, err := account.FromSeed(u.Wallet)
	if err != nil {
		panic(err)
	}
	if !account.IsValidAccountNumber(acct.AccountNumber) {
		panic("Account is not valid")
	}
	return acct
}

type Message struct {
	Title      string
	Content    string
	Target     string
	TargetName string
}
