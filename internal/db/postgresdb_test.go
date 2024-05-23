package db_test

import (
	"context"
	"testing"

	"postsandcomments/graph/model"
	"postsandcomments/internal/db"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *db.PostgresDB {
	db, err := db.NewPostgresDB("db", 5432, "postgres", "password")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func TestNewPostgresDB(t *testing.T) {
	db, err := db.NewPostgresDB("db", 5432, "postgres", "password")
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestCreatePostPostgres(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	post := &model.Post{
		ID:            uuid.New().String(),
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}

	err := db.CreatePost(context.Background(), post)
	assert.NoError(t, err)

	posts, err := db.GetPosts(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, post, posts[0])
}

func TestGetPostsPostgres(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	post1 := &model.Post{
		ID:            uuid.New().String(),
		Title:         "Post 1",
		Body:          "Body 1",
		AllowComments: true,
	}
	post2 := &model.Post{
		ID:            uuid.New().String(),
		Title:         "Post 2",
		Body:          "Body 2",
		AllowComments: false,
	}

	err := db.CreatePost(context.Background(), post1)
	assert.NoError(t, err)
	err = db.CreatePost(context.Background(), post2)
	assert.NoError(t, err)

	posts, err := db.GetPosts(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 2, len(posts))
	assert.Equal(t, posts[0], post1)
	assert.Equal(t, posts[1], post2)

}

func TestGetPostByIdPostgres(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	post := &model.Post{
		ID:            uuid.New().String(),
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}

	err := db.CreatePost(context.Background(), post)
	assert.NoError(t, err)

	fetchedPost, err := db.GetPostById(context.Background(), post.ID, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, post, fetchedPost)
}

func TestGetCommentByIdPostgres(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()


	post := &model.Post{
		ID:            uuid.New().String(),
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}
	comment := &model.Comment{
		ID:       uuid.New().String(),
		PostID:   post.ID,
		Body:     "Test Comment",
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

func TestCreateCommentPostgres(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	post := &model.Post{
		ID:            uuid.New().String(),
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}

	err := db.CreatePost(context.Background(), post)
	assert.NoError(t, err)

	comment := &model.Comment{
		ID:       uuid.New().String(),
		PostID:   post.ID,
		Body:     "Test Comment",
		ParentID: nil,
	}

	err = db.CreateComment(context.Background(), post, comment)
	assert.NoError(t, err)

	childrenComment := &model.Comment{
		ID:       uuid.New().String(),
		PostID:   post.ID,
		Body:     "Test Comment",
		ParentID: &comment.ID,
	}
	db.CreateComment(context.Background(), post, childrenComment)
	post, err = db.GetPostById(context.Background(), post.ID, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, childrenComment, post.Comments[1])
}

func TestCreateCommentWithLimitAndOffsetPostgres(t *testing.T) {
	db := setupTestDB(t)
	defer db.DB.Close()

	post := &model.Post{
		ID:            uuid.New().String(),
		Title:         "Test Post",
		Body:          "Test body",
		AllowComments: true,
	}

	err := db.CreatePost(context.Background(), post)
	assert.NoError(t, err)

	comments := []*model.Comment{
		{
			ID:       uuid.New().String(),
			PostID:   post.ID,
			Body:     "Test Comment 1",
			ParentID: nil,
		},
		{
			ID:       uuid.New().String(),
			PostID:   post.ID,
			Body:     "Test Comment 2",
			ParentID: nil,
		},
	}
	childComment := &model.Comment{
		ID:       uuid.New().String(),
		PostID:   post.ID,
		Body:     "Test Comment 1-1",
		ParentID: &comments[0].ID,
	}
	comments = append(comments, childComment)

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
