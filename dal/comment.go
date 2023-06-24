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
	var user = models.User{ID: userID}
	dbc := db.First(&user)
	if dbc.Error != nil && dbc.Error.Error() == "record not found" {
		newUser := &models.User{ID: userID}
		dbc = db.WithContext(*ctx).Create(newUser)
		if dbc.Error != nil {
			return nil, dbc.Error
		}
	}

	newComment := &models.Comment{
		UserID:    userID,
		ReplyToId: replyToId,
		Content:   content,
		Agent:     agent,
	}
	newComment.PopulateIdentifier(outerID)
	dbc = db.WithContext(*ctx).Create(newComment)
	return newComment, dbc.Error
}
