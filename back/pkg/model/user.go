package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserID     string    `gorm:"column:user_id;primary_key" json:"userId,omitempty"`
	Status     string    `gorm:"column:status" json:"status,omitempty"`
	UserName   string    `gorm:"column:user_name;not null;unique" json:"userName,omitempty"`
	Password   string    `gorm:"column:password" json:"password,omitempty"`
	Token      string    `gorm:"column:token" json:"token,omitempty"`
	Address    string    `gorm:"column:address" json:"address,omitempty"`
	Role       string    `gorm:"column:role" json:"role,omitempty"`
	LoginTime  time.Time `gorm:"column:login_time" json:"omitempty"`
	LockedTime time.Time `gorm:"column:locked_time" json:"locked_time,omitempty"`
	ExpireTime time.Time `gorm:"column:expire_time" json:"expire_time,omitempty"`
	Label      string    `gorm:"column:label" json:"label,omitempty"`
}
