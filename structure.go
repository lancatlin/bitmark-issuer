package main

import (
	"time"

	"github.com/bitmark-inc/bitmark-sdk-go/account"
)

// Asset record a file an its every issues
type Asset struct {
	// A random id
	ID        string `gorm:"primary_key"`
	AssetID   string
	Name      string `form:"asset_name"`
	Amount    int    `form:"amount"`
	CreatedAt time.Time
	Issues    []Issue
	User      User
	UserID    uint
	URL       *URL
}

// Issue belongs to an Asset
// Issue record its Reciever
type Issue struct {
	ID      string
	Asset   Asset
	AssetID string
	// 如果尚未轉移留空
	Receiver string
	// 如果尚未轉移留空
	TransferredAt time.Time
	CreatedAt     time.Time
	TxID          string
}

// User ...
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

// Account return the bitmark account of a user
func (u User) Account() account.Account {
	acct, err := account.FromSeed(u.Wallet)
	if err != nil {
		panic(err)
	}
	/*
		This function is not use anymore
		if !account.IsValidAccountNumber(acct.AccountNumber) {
			panic("Account is not valid")
		}
	*/
	return acct
}

// URL is the url token to access a asset
type URL struct {
	ID        string `gorm:"primary_key"`
	Asset     Asset
	AssetID   string
	CreatedAt time.Time
	ExpireAt  time.Time
}

func (u URL) String() string {
	return u.ID
}

type message struct {
	Title      string
	Content    string
	Target     string
	TargetName string
}
