package db

import (
	"context"
	"postsandcomments/internal/graph/model"
)

type Database interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPosts(ctx context.Context) ([]*model.Post, error)
	GetPostById(ctx context.Context, id string, limit *int, offset *int) (*model.Post, error)
	CreateComment(ctx context.Context, post *model.Post, comment *model.Comment) error
	GetCommentById(ctx context.Context, id string) (*model.Comment, error)
}
