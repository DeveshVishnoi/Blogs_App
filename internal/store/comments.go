package store

import (
	"context"
	"database/sql"
)

type CommentStore struct {
	db *sql.DB
}

type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

func (cm *CommentStore) GetByPostId(ctx context.Context, postId int) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id  FROM comments c
		JOIN users on users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;
	`

	rows, err := cm.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}

		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.UserName, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (cm *CommentStore) Create(ctx context.Context, comment *Comment) error {

	query := `
	
	insert into comments (post_id, user_id, content)
	values ($1, $2,$3)
	returning id, created_at
	`

	err := cm.db.QueryRowContext(ctx, query,
		comment.PostID,
		comment.UserID,
		comment.Content).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return err
	}

	return nil

}
