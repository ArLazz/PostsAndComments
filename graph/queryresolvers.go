package graph

import (
	"context"
	"fmt"
	"postsandcomments/graph/model"
)

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.DataBase.GetPosts(ctx)
	if err != nil {
		r.Logger.Errorf("error to get all posts: %v", err)
		return nil, fmt.Errorf("error to get all posts: %v", err)
	}

	r.Logger.Infof("get all posts")
	return posts, nil
}

func (r *queryResolver) Post(ctx context.Context, id string, limit *int, offset *int) (*model.Post, error) {
	post, err := r.DataBase.GetPostById(ctx, id, limit, offset)
	if err != nil {
		r.Logger.Errorf("error to get post by id: %v", err)
		return nil, fmt.Errorf("error to get post by id: %v", err)
	}

	r.Logger.Infof("get post with id = %s", id)
	return post, nil
}
