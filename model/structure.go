package model

import (
	"github.com/jinzhu/gorm"
	"time"
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
	gorm.Model
	Email string `gorm:"unique"`
	// 簽署用名稱
	Name     string
	Password []byte
	// 儲存 account seed
	Wallet string
	Assets []Asset
}
