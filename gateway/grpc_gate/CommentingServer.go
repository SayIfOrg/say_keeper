package grpc_gate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SayIfOrg/say_keeper/commenting"
	"github.com/SayIfOrg/say_keeper/dal"
	"github.com/SayIfOrg/say_keeper/graph/gmodel"
	"github.com/SayIfOrg/say_keeper/models"
	pb "github.com/SayIfOrg/say_protos/packages/go"
	"github.com/redis/go-redis/v9"
	"gopkg.in/guregu/null.v4/zero"
	"gorm.io/gorm"
	"log"
)

// CommentingServer is used to implement webpage.ManageInstanceServer.
type CommentingServer struct {
	pb.UnsafeCommentingServer
	DB  *gorm.DB
	RDB *redis.Client
}

func (s *CommentingServer) Post(ctx context.Context, in *pb.PostComment) (*pb.Comment, error) {
	var replyToComment models.Comment

	if in.GetReplyToOuterIdentifier() != "" {
		var platform string
		var agent = "telebot" // TODO
		if agent == commenting.TelegramAgent {
			platform = "telebot_id"
		} else {
			return nil, errors.New("agent not supported")
		}
		dbc := s.DB.Joins(
			"left join comment_plats on comment_plats.comment_id = comments.id").Where(
			fmt.Sprintf("comment_plats.%s = ?", platform), in.GetReplyToOuterIdentifier()).First(&replyToComment)
		if dbc.Error != nil {
			return nil, dbc.Error
		}
	}

	comment, err := dal.CreateComment(
		&ctx,
		s.DB,
		uint(in.GetUserId()),
		zero.IntFrom(int64(replyToComment.ID)).NullInt64,
		in.Content,
		commenting.TelegramAgent,
		in.GetOuterIdentifier(),
	)
	if err != nil {
		return nil, err
	}

	gComment := gmodel.FromDBComment(comment)
	jsonBytes, err := json.Marshal(gComment)
	if err != nil {
		panic(err)
	}
	err = commenting.PublishComment(s.RDB, ctx, jsonBytes)
	if err != nil {
		// TODO proper handle this case
		log.Println("error in publishing the comment", err)
	}
	return &pb.Comment{
		Id:              uint64(comment.ID),
		UserId:          uint64(comment.UserID),
		ReplyToId:       uint64(comment.ReplyToId.Int64),
		Content:         comment.Content,
		OuterIdentifier: comment.Platform.TelebotID.String,
	}, nil
}

func (s *CommentingServer) Get(ctx context.Context, in *pb.GetParams) (*pb.Comment, error) {
	var comment models.Comment
	dbc := s.DB.Joins(
		"left join comment_plats on comment_plats.comment_id = comments.id").Where(
		"comment_plats.telebot_id = ?", in.GetOuterIdentifier()).First(&comment)
	if dbc.Error != nil {
		if dbc.Error.Error() == "record not found" {
			return &pb.Comment{}, nil
		}
		return nil, dbc.Error
	}

	return &pb.Comment{
		Id:              uint64(comment.ID),
		OuterIdentifier: comment.Platform.TelebotID.String,
		UserId:          uint64(comment.UserID),
		ReplyToId:       uint64(comment.ReplyToId.Int64),
		Content:         comment.Content,
	}, nil
}
