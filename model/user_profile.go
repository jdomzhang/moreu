package model

import "time"

type UserProfile struct {
	Id       int64     `json:"id"`
	UserId   int64     `json:"user_id" gorm:"not null"`
	Nickname string    `json:"nickname" gorm:"size:32;not null"`
	Avatar   string    `json:"avatar" gorm:"size:255;not null"`
	Bio      string    `json:"bio" gorm:"size:255;not null"`
	URL      string    `json:"url" gorm:"size:255;not null"`
	Company  string    `json:"company" gorm:"size:32;not null"`
	Location string    `json:"location" gorm:"size:32;not null"`
	Deleted  time.Time `json:"deleted" gorm:"column:deleted_at;not null"`
	Created  time.Time `json:"created" gorm:"column:created_at;not null"`
	Updated  time.Time `json:"updated" gorm:"column:updated_at;not null"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}