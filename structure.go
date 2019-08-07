package main

import (
	"time"
	"io"
)

type Asset struct {
	// bitmark asset id
	ID       string `gorm:"primary_key"`
	Amount   int
	FilePath string
	file     io.Reader 
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
	ID uint 
	CreatedAt time.Time 
	UpdatedAt time.Time
	Email string `gorm:"unique" form:"email"`
	// 簽署用名稱
	Name     string `form:"name"`
	PlainPassword string `form:"password" gorm:"-"`
	Password []byte
	// 儲存 account seed
	Wallet string
	Assets []Asset
	IsLogin bool `gorm:"-"`
}

type Message struct {
	Title string
	Content string
	Target string
	TargetName string
}