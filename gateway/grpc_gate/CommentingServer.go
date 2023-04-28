package grpc_gate

import (
	"context"
	"encoding/json"
	"github.com/SayIfOrg/say_keeper/commenting"
	"github.com/SayIfOrg/say_keeper/dal"
	"github.com/SayIfOrg/say_keeper/graph/gmodel"
	pb "github.com/SayIfOrg/say_protos/packages/go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
)

// CommentingServer is used to implement webpage.ManageInstanceServer.
type CommentingServer struct {
	pb.UnsafeCommentingServer
	DB  *gorm.DB
	RDB *redis.Client
}

func (s *CommentingServer) Post(ctx context.Context, in *pb.Comment) (*pb.Comment, error) {
	var replyToID uint
	var pReplyToID *uint
	if in.GetReplyToId() != 0 {
		replyToID = uint(in.GetReplyToId())
		pReplyToID = &replyToID
	}
	comment, err := dal.CreateComment(
		&ctx,
		s.DB,
		uint(in.GetUserId()),
		pReplyToID,
		in.Content,
		commenting.TelegramAgent,
		in.GetOuterIdentifier(),
	)
	if err != nil {
		return nil, err
	}
	var dReplyToID uint64
	if comment.ReplyToId != nil {
		dReplyToID = uint64(*comment.ReplyToId)
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
		ReplyToId:       dReplyToID,
		Content:         comment.Content,
		OuterIdentifier: comment.Identifier.String,
	}, nil
}
