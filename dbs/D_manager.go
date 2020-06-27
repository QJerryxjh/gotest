package dbs

import (
	"time"
)

type DbManger struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json: "username" gorm:"type:varchar(16);not null;unique_index"`
	Password  string    `json: "password" gorm:"type:varchar(32);not null"`
	Name      string    `json: "name" gorm:"type:varchar(64);not null;"`
	Email     string    `json: "email" gorm:"type:varchar(32)"`
	Gender    string    `json: "gender" gorm:"type:varchar(16)"`
	Birthday  time.Time `json: "birthday" gorm:"type:datetime;default:null"`
	CreatedAt time.Time `json: "created_at" gorm:"type:datetime"`
}

func (self *DbManger) QueryAllManager() (manager DbManger, err error) {
	db := DB

	db = db.First(&manager)
	return
}
