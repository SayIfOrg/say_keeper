package models

import "database/sql"

type User struct {
	ID       uint      `gorm:"primaryKey;autoIncrement:false"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}

type Comment struct {
	ID         uint
	Identifier sql.NullString `gorm:"unique"`
	UserID     uint           `gorm:"not null"`
	ReplyToId  sql.NullInt64
	Replies    []Comment `gorm:"foreignKey:ReplyToId"`
	Content    string    `gorm:"not null"`
	Agent      string    `gorm:"not null"`
}
