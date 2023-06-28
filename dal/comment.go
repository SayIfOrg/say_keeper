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
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newComment).Error; err != nil {
			return err
		}
		commentPlat := &models.CommentPlat{
			CommentID: newComment.ID,
			TelebotID: sql.NullString{
				String: outerID,
				Valid:  true,
			},
		}
		if err := tx.Create(commentPlat).Error; err != nil {
			return err
		}

		return nil
	})

	return nil, err
}
