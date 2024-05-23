package db

import (
	"context"
	"fmt"
	"postsandcomments/graph/model"
	"sync"
)

type InMemoryDB struct {
	Posts    map[string]*model.Post
	Comments map[string]*model.Comment
	Mutex    sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Posts:    make(map[string]*model.Post),
		Comments: make(map[string]*model.Comment),
		Mutex:    sync.RWMutex{},
	}
}

func (db *InMemoryDB) CreatePost(ctx context.Context, post *model.Post) error {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	db.Posts[post.ID] = post

	return nil
}

func (db *InMemoryDB) GetPosts(ctx context.Context) ([]*model.Post, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	posts := make([]*model.Post, 0, len(db.Posts))
	for _, post := range db.Posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (db *InMemoryDB) GetCommentById(ctx context.Context, id string) (*model.Comment, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	comment, exists := db.Comments[id]
	if !exists {
		return nil, fmt.Errorf("no comments with this id: %s", id)
	}

	return comment, nil
}

func (db *InMemoryDB) GetPostById(ctx context.Context, id string, limit *int, offset *int) (*model.Post, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	post, exists := db.Posts[id]
	if !exists {
		return nil, fmt.Errorf("no posts with this id: %s", id)
	}

	if limit != nil && offset != nil {
		return &model.Post{
			ID:            post.ID,
			Title:         post.Title,
			Body:          post.Body,
			AllowComments: post.AllowComments,
			Comments:      GetComments(post.Comments, *limit, *offset),
		}, nil
	}

	return post, nil
}
func GetComments(comments []*model.Comment, limit int, offset int) []*model.Comment {
	result := make([]*model.Comment, 0)
	startIndex, endIndex := offset, limit+offset
	var processComents func(comments []*model.Comment, startIndex, endIndex int)
	processComents = func(comments []*model.Comment, startIndex, endIndex int) {
		if len(result) == endIndex {
			return
		}
		for _, comment := range comments {

			if len(result) < endIndex {
				result = append(result, comment)
			}
		}

		for _, comment := range comments {
			processComents(comment.Children, startIndex, endIndex)
		}
	}

	processComents(comments, startIndex, endIndex)
	if endIndex > len(result) {
		result = result[startIndex:]
	} else {
		result = result[startIndex:endIndex]
	}

	return result
}

func (db *InMemoryDB) CreateComment(ctx context.Context, post *model.Post, comment *model.Comment) error {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	if comment.ParentID != nil {
		ParentComment, exists := db.Comments[*comment.ParentID]
		if !exists {
			return fmt.Errorf("no comments with this id: %s", *comment.ParentID)
		}

		ParentComment.Children = append(ParentComment.Children, comment)
	} else {
		db.Comments[comment.ID] = comment
		post.Comments = append(post.Comments, comment)
	}
	return nil
}
