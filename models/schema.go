package models

import "database/sql"

type User struct {
	ID       uint
	Name     string
	Comments []Comment `gorm:"foreignKey:UserID"`
}

type Comment struct {
	ID         uint
	Identifier sql.NullString `gorm:"unique"`
	UserID     uint
	ReplyToId  *uint
	Replies    []Comment `gorm:"foreignKey:ReplyToId"`
	Content    string
	Agent      string
}
