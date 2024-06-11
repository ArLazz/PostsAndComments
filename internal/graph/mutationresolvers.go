package graph

import (
	"context"
	"fmt"
	"postsandcomments/internal/graph/model"
	"unicode/utf8"

	"github.com/google/uuid"
)

const MaxLengthOfComment = 2000

func (r *mutationResolver) CreatePost(ctx context.Context, title string, body string, allowComments bool) (*model.Post, error) {
	post := &model.Post{
		ID:            uuid.New().String(),
		Title:         title,
		Body:          body,
		Comments:      make([]*model.Comment, 0),
		AllowComments: allowComments,
	}

	err := r.DataBase.CreatePost(ctx, post)
	if err != nil {
		r.Logger.Errorf("error to create post: %v", err)
		return nil, fmt.Errorf("error to create post: %v", err)
	}

	r.Logger.Infof("post with id = %s created", post.ID)
	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, postID string, body string, parentID *string) (*model.Comment, error) {
	comment := &model.Comment{
		ID:       uuid.New().String(),
		PostID:   postID,
		Body:     body,
		ParentID: parentID,
		Children: make([]*model.Comment, 0),
	}

	if utf8.RuneCountInString(body) > MaxLengthOfComment {
		r.Logger.Errorf("error to create comment: size of comment more than max size")
		return nil, fmt.Errorf("error to create comment: size of comment more than max size")
	}

	post, err := r.DataBase.GetPostById(ctx, postID, nil, nil)
	if err != nil {
		r.Logger.Errorf("error to get post by id to create comment: %v", err)
		return nil, fmt.Errorf("error to get post by id to create comment: %v", err)
	}

	if comment.ParentID != nil {
		parentComment, err := r.DataBase.GetCommentById(ctx, *comment.ParentID)
		if err != nil {
			r.Logger.Errorf("error to get parent comment by id to create comment: %v", err)
			return nil, fmt.Errorf("error to get parent comment by id to create comment: %v", err)
		}
		if parentComment.PostID != comment.PostID {
			r.Logger.Errorf("postID for parent and child comment should be the same")
			return nil, fmt.Errorf("postID for parent and child comment should be the same")
		}
	}

	if !post.AllowComments {
		r.Logger.Errorf("error to create comment: not allowed comments for post")
		return nil, fmt.Errorf("error to create comment: not allowed comments for post")
	}

	err = r.DataBase.CreateComment(ctx, post, comment)
	if err != nil {
		r.Logger.Errorf("error to create comment: %v", err)
		return nil, fmt.Errorf("error to create comment: %v", err)
	}

	event := &CommentEvent{
		PostID:  postID,
		Comment: comment,
	}
	r.SubscriptionManager.Publish(event)

	r.Logger.Infof("comment with id = %s created", comment.ID)
	return comment, err
}
