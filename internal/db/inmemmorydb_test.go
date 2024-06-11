package db_test

import (
	"context"
	"postsandcomments/internal/graph/model"
	"postsandcomments/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryDB(t *testing.T) {
	db := db.NewInMemoryDB()

	assert.NotNil(t, db)
	assert.Equal(t, make(map[string]*model.Post), db.Posts)
	assert.Equal(t, make(map[string]*model.Comment), db.Comments)
}

func TestCreatePostInMemory(t *testing.T) {
	db := db.NewInMemoryDB()

	post := &model.Post{
		ID:            "test_post_id",
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}
	err := db.CreatePost(context.Background(), post)
	assert.NoError(t, err)
	assert.Contains(t, db.Posts, post.ID)
}

func TestGetPostsInMemory(t *testing.T) {
	db := db.NewInMemoryDB()

	post1 := &model.Post{
		ID:            "test_post_1",
		Title:         "First Post",
		Body:          "Test body",
		AllowComments: true,
	}
	post2 := &model.Post{
		ID:            "test_post_2",
		Title:         "Second Post",
		Body:          "Test body",
		AllowComments: false,
	}

	db.CreatePost(context.Background(), post1)
	db.CreatePost(context.Background(), post2)

	posts, _ := db.GetPosts(context.Background())
	assert.Len(t, posts, 2)
	assert.Contains(t, posts, post1)
	assert.Contains(t, posts, post2)
}

func TestGetPostByIdInMemory(t *testing.T) {
	db := db.NewInMemoryDB()

	post := &model.Post{
		ID:            "test_post_id",
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}
	
	err := db.CreatePost(context.Background(), post)
	assert.NoError(t, err)

	retrievedPost, err := db.GetPostById(context.Background(), post.ID, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, post, retrievedPost)
}
func TestGetCommentByIdInMemory(t *testing.T) {
	db := db.NewInMemoryDB()

	post := &model.Post{
		ID:            "test_post_id",
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}
	comment := &model.Comment{
		ID:       "comment_post_1",
		PostID:   post.ID,
		Body:     "Test Comment 1",
		ParentID: nil,
	}
	err := db.CreatePost(context.Background(), post)
	assert.NoError(t, err)
	err = db.CreateComment(context.Background(), post, comment)
	assert.NoError(t, err)

	retrievedComment, err := db.GetCommentById(context.Background(), comment.ID)
	assert.NoError(t, err)
	assert.Equal(t, comment, retrievedComment)
}

func TestGetPostByIdWithLimitAndOffsetInMemory(t *testing.T) {
	db := db.NewInMemoryDB()

	post := &model.Post{
		ID:            "test_post_id",
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}
	err := db.CreatePost(context.Background(), post)
	assert.NoError(t, err)

	comments := []*model.Comment{
		{
			ID:       "comment_post_1",
			PostID:   post.ID,
			Body:     "Test Comment 1",
			ParentID: nil,
		},
		{
			ID:       "comment_post_2",
			PostID:   post.ID,
			Body:     "Test Comment 2",
			ParentID: nil,
		},
	}
	comments = append(comments, &model.Comment{
		ID:       "comment_post_1-1",
		PostID:   post.ID,
		Body:     "Test Comment 1-1",
		ParentID: &comments[0].ID,
	})

	for _, comment := range comments {
		err = db.CreateComment(context.Background(), post, comment)
		assert.NoError(t, err)
	}

	limit, offset := 1, 1
	fetchedPost, err := db.GetPostById(context.Background(), post.ID, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, comments[1], fetchedPost.Comments[0])
	limit, offset = 1, 2
	fetchedPost, err = db.GetPostById(context.Background(), post.ID, &limit, &offset)
	assert.NoError(t, err)
	assert.Equal(t, comments[2], fetchedPost.Comments[0])
}

func TestCreateCommentInMemory(t *testing.T) {
	db := db.NewInMemoryDB()

	post := &model.Post{
		ID:            "test_post_id",
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}
	comment := &model.Comment{
		ID:       "comment_post_1",
		PostID:   post.ID,
		Body:     "Test Comment 1",
		ParentID: nil,
	}
	db.CreatePost(context.Background(), post)

	err := db.CreateComment(context.Background(), post, comment)
	assert.NoError(t, err)
	assert.Contains(t, post.Comments, comment)
	assert.Contains(t, db.Comments, comment.ID)

	childrenComment := &model.Comment{
		ID:       "child_comment_id",
		PostID:   post.ID,
		Body:     "Test Comment child",
		ParentID: &comment.ID,
	}
	db.CreateComment(context.Background(), post, childrenComment)
	assert.Contains(t, comment.Children, childrenComment)
}
