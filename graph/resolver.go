package graph

//go:generate go run github.com/99designs/gqlgen generate
// go generate ./...

import (
	"github.com/SayIfOrg/say_keeper/commenting"
	"github.com/SayIfOrg/say_keeper/graph/gmodel"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB   *gorm.DB
	RDB  *redis.Client
	Subs *commenting.Subs[gmodel.Comment]
}
