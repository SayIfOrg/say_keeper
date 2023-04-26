package dal

import (
	"context"
	"github.com/SayIfOrg/say_keeper/models"
	"gorm.io/gorm"
)

func CreateComment(
	ctx *context.Context,
	db *gorm.DB,
	userID uint,
	replyToId *uint,
	content string,
	agent string,
) (*models.Comment, error) {
	newComment := &models.Comment{
		UserID:    userID,
		ReplyToId: replyToId,
		Content:   content,
		Agent:     agent,
	}
	db.WithContext(*ctx).Create(newComment)
	return newComment, db.Error
}
