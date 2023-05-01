package dal

import (
	"context"
	"database/sql"
	"github.com/SayIfOrg/say_keeper/models"
	"gorm.io/gorm"
)

func CreateComment(
	ctx *context.Context,
	db *gorm.DB,
	userID uint,
	replyToId sql.NullInt64,
	content string,
	agent string,
	outerID string,
) (*models.Comment, error) {
	newComment := &models.Comment{
		UserID:    userID,
		ReplyToId: replyToId,
		Content:   content,
		Agent:     agent,
	}
	newComment.PopulateIdentifier(outerID)
	dbc := db.WithContext(*ctx).Create(newComment)
	return newComment, dbc.Error
}
