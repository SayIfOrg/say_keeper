package dataloader

import (
	"context"
	"fmt"
	"github.com/SayIfOrg/say_keeper/graph/gmodel"
	"github.com/SayIfOrg/say_keeper/models"
	"github.com/SayIfOrg/say_keeper/utils"
	"github.com/graph-gophers/dataloader"
	"gorm.io/gorm"
	"strconv"
)

// GetComment wraps the Comment dataloader for efficient retrieval by user ID
func (i *Loaders) GetComment(ctx context.Context, commentID string) (*gmodel.Comment, error) {
	thunk := i.CommentLoader.Load(ctx, dataloader.StringKey(commentID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*gmodel.Comment), nil
}

// CommentReader used as dependency injection for reading comments
type CommentReader struct {
	conn *gorm.DB
}

// LoadComments implements a batch function that can retrieve many comments by ID,
func (u *CommentReader) LoadComments(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// collect the keys to search for
	var commentIDs = make([]string, len(keys))
	for ix, key := range keys {
		commentIDs[ix] = key.String()
	}
	// read all requested comments in a single query
	var comments []models.Comment
	dbc := u.conn.WithContext(ctx).Find(&comments, commentIDs)
	if dbc.Error != nil {
		panic(dbc.Error)
	}
	commentById := map[string]*gmodel.Comment{}
	for _, comment := range comments {
		commentById[strconv.Itoa(int(comment.ID))] = &gmodel.Comment{
			ID:        strconv.Itoa(int(comment.ID)),
			UserID:    strconv.Itoa(int(comment.UserID)),
			ReplyToID: utils.RUintToString(comment.ReplyToId),
			Content:   comment.Content,
			Agent:     comment.Agent,
		}
	}

	// return comments in the same order requested
	output := make([]*dataloader.Result, len(keys))
	for index, commentKey := range keys {
		comment, ok := commentById[commentKey.String()]
		if ok {
			output[index] = &dataloader.Result{Data: comment, Error: nil}
		} else {
			err := fmt.Errorf("user not found %s", commentKey.String())
			output[index] = &dataloader.Result{Data: nil, Error: err}
		}
	}
	return output
}
