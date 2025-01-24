package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Post struct {
	ID        int       `json:"id" bson:"id"`
	Content   string    `json:"content" bson:"content"`
	Title     string    `json:"title" bson:"title"`
	UserID    int       `json:"user_id" bson:"user_id"`
	Tags      []string  `json:"tags" bson:"tags"`
	CreatedAt string    `json:"created_at" bson:"created_at"`
	UpdatedAt string    `json:"updated_at" bson:"updated_at"`
	Comments  []Comment `json:"comment"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db *sql.DB
}

func (ps *PostStore) Create(ctx context.Context, post *Post) error {

	query := `
	Insert into posts (content, title, user_id, tags)
	Values ($1, $2, $3, $4) returning  id, created_At, updated_at
	`

	err := ps.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostStore) GetById(ctx context.Context, postId int) (*Post, error) {
	query := `

	select id, user_id, title, content, created_At,  updated_at,tags
	from posts
	where id = $1

	`

	var post Post

	err := ps.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, pq.Array(&post.Tags),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	fmt.Println("Post -- ", post)
	return &post, nil
}

func (ps *PostStore) Delete(ctx context.Context, postId int) error {
	query := `
	delete from posts where id = $1
	`

	res, err := ps.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (ps *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
	update posts set title = $1 , content = $2 where id = $3
	`

	_, err := ps.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostStore) GetUserFeed(ctx context.Context, userId int, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {

	query := `
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.created_at, 
			p.tags, 
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c 
			ON c.post_id = p.id
		LEFT JOIN users u 
			ON p.user_id = u.id
		WHERE 
			p.user_id = $1 -- Include posts authored by the user
			OR EXISTS (
				SELECT $1
				FROM followers f
				WHERE f.user_id = $1
				AND f.follower_id = p.user_id
			)
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + ` 
		Limit $2 Offset $3
	`

	rows, err := ps.db.QueryContext(ctx, query, userId, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fmt.Println("rows -- ", rows)
	var feed []PostWithMetaData
	for rows.Next() {
		var p PostWithMetaData
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			pq.Array(&p.Tags),
			&p.User.UserName,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		fmt.Println("P -", p)

		feed = append(feed, p)
	}

	return feed, nil
}
