package model

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "data.sqlite")
	if err != nil {
		panic(err)
	}
	db.AutoMigration(&Asset{}, &Issue{})
}
