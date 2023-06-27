package gmodel

import "time"

type Comment struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userID"`
	ReplyToID *string    `json:"replyToID"`
	ReplyTo   *Comment   `json:"replyTo"`
	Replies   []*Comment `json:"replies"`
	Content   string     `json:"content"`
	Agent     string     `json:"agent"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
