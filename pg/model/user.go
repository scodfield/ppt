package model

import "gorm.io/gorm"

type User struct {
	UserID   int64  `gorm:"primaryKey" json:"user_id"`
	Username string `gorm:"size:255" json:"username"`
	Password string `gorm:"size:255" json:"password"`
	Email    string `gorm:"uniqueIndex;size:255" json:"email"`
}

func (User) TableName() string {
	return "user"
}

func MigrateUserModel(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
