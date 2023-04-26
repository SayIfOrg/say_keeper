package dataloader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"gorm.io/gorm"
	"net/http"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	//Add the other loaders here
	CommentLoader *dataloader.Loader
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *gorm.DB) *Loaders {
	// define the data loader
	commentReader := &CommentReader{conn: conn}
	loaders := &Loaders{
		CommentLoader: dataloader.NewBatchedLoader(commentReader.LoadComments),
	}
	return loaders
}

// Middleware injects data loaders into the context
func Middleware(loaders *Loaders, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
