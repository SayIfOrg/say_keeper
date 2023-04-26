package grpc_gate

import (
	"context"
	pb "github.com/SayIfOrg/say_protos/packages/go"
)

// CommentingServer is used to implement webpage.ManageInstanceServer.
type CommentingServer struct {
	pb.UnsafeCommentingServer
}

func (s *CommentingServer) Post(ctx context.Context, in *pb.Comment) (*pb.Comment, error) {
	return nil, nil
}
