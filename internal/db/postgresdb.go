package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"postsandcomments/graph/model"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(host string, port int, user, password string) (*PostgresDB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)

	db, err := connectToDB(psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error to create postgres db: %v", err)
	}
	
	_, err = db.Exec(`
		DROP TABLE IF EXISTS comments;
		DROP TABLE IF EXISTS posts;

		CREATE TABLE posts (
			id UUID PRIMARY KEY,
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			allow_comments BOOLEAN NOT NULL
		);

		CREATE TABLE comments (
			id UUID PRIMARY KEY,
			post_id UUID REFERENCES posts(id),
			body VARCHAR(2000) NOT NULL,
			parent_id UUID,
			FOREIGN KEY (parent_id) REFERENCES comments (id)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("error to create tables: %v", err)
	}

	return &PostgresDB{DB: db}, nil
}

func (db *PostgresDB) CreatePost(ctx context.Context, post *model.Post) error {
	query := `INSERT INTO posts (id, title, body, allow_comments) VALUES ($1, $2, $3, $4)`
	_, err := db.DB.ExecContext(ctx, query, post.ID, post.Title, post.Body, post.AllowComments)
	return err
}

func (db *PostgresDB) GetPosts(ctx context.Context) ([]*model.Post, error) {
	rows, err := db.DB.QueryContext(ctx, "SELECT id, title, body, allow_comments FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Body,
			&post.AllowComments,
		)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (db *PostgresDB) GetPostById(ctx context.Context, id string, limit *int, offset *int) (*model.Post, error) {
	row := db.DB.QueryRowContext(ctx, "SELECT id, title, body, allow_comments FROM posts WHERE id=$1", id)
	var post model.Post

	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Body,
		&post.AllowComments,
	)
	if err != nil {
		return nil, err
	}

	comments, err := db.GetComments(ctx, id, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error to get comments: %v", err)
	}
	post.Comments = comments

	return &post, nil
}

func (db *PostgresDB) GetCommentById(ctx context.Context, id string) (*model.Comment, error) {
	row := db.DB.QueryRowContext(ctx, "SELECT id, post_id, body, parent_id FROM comments WHERE id=$1", id)
	var comment model.Comment

	err := row.Scan(
		&comment.ID,
		&comment.PostID,
		&comment.Body,
		&comment.ParentID,
	)
	if err != nil {
		return nil, err
	}
	
	return &comment, nil
}

func (db *PostgresDB) CreateComment(ctx context.Context, post *model.Post, comment *model.Comment) error {
	query := `INSERT INTO comments (id, post_id, body, parent_id) VALUES ($1, $2, $3, $4)`
	_, err := db.DB.ExecContext(ctx, query, comment.ID, post.ID, comment.Body, comment.ParentID)
	return err
}

func (db *PostgresDB) GetComments(ctx context.Context, postId string, limit *int, offset *int) ([]*model.Comment, error) {
	query := `
        WITH RECURSIVE comment_tree AS (
            SELECT
                id,
                post_id,
                body,
                parent_id
            FROM comments
            WHERE post_id =  $1 AND parent_id IS NULL

            UNION ALL

            SELECT
                c.id,
                c.post_id,
                c.body,
                c.parent_id
            FROM comments c
            INNER JOIN comment_tree ct ON c.parent_id = ct.id
        )
        SELECT id, post_id, body, parent_id FROM comment_tree LIMIT $2 OFFSET $3;
    `
	rows, err := db.DB.QueryContext(ctx, query, postId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment

	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Body, &comment.ParentID); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func connectToDB(psqlInfo string) (*sql.DB, error) {
	var db *sql.DB
    var err error
	for i := 1; i < 11; i++{
		logrus.Infof("waiting for inicialization of db, attempt %d", i)
		time.Sleep(time.Second)
        db, err = sql.Open("postgres", psqlInfo)
        if err != nil {
            logrus.Errorf("failed to open database connection: %v", err)
            continue
        }

        err = db.Ping()
        if err != nil {
            logrus.Errorf("ping failed: %v", err)
            continue
        }

        return db, nil 
    }
	return nil, fmt.Errorf("error to connect to db")
}