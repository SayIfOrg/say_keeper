package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID       uint      `gorm:"primaryKey;autoIncrement:false"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}

type CommentPlat struct {
	CommentID uint
	TelebotID sql.NullString
}

type Comment struct {
	ID        uint
	Platform  CommentPlat `gorm:"foreignKey:CommentID"`
	UserID    uint        `gorm:"not null"`
	ReplyToId sql.NullInt64
	Replies   []Comment `gorm:"foreignKey:ReplyToId"`
	Content   string    `gorm:"not null"`
	Agent     string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
