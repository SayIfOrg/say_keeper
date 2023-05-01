package grpc_gate

import (
	"context"
	"encoding/json"
	"github.com/SayIfOrg/say_keeper/commenting"
	"github.com/SayIfOrg/say_keeper/dal"
	"github.com/SayIfOrg/say_keeper/graph/gmodel"
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

func (s *CommentingServer) Post(ctx context.Context, in *pb.Comment) (*pb.Comment, error) {
	comment, err := dal.CreateComment(
		&ctx,
		s.DB,
		uint(in.GetUserId()),
		zero.IntFrom(int64(in.GetReplyToId())).NullInt64,
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
		OuterIdentifier: comment.Identifier.String,
	}, nil
}
