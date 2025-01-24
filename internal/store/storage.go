package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("record Not Found")
	ErrConflict          = errors.New("resource already exists")
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int) (*Post, error)
		Delete(context.Context, int) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int, PaginatedFeedQuery) ([]PostWithMetaData, error)
	}

	Users interface {
		Create(context.Context, *User, *sql.Tx) error
		GetById(context.Context, int) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		createUserInvitation(context.Context, *sql.Tx, string, time.Duration, int) error
		Activate(context.Context, string) error
		Delete(context.Context, int) error
	}

	Comments interface {
		GetByPostId(context.Context, int) ([]Comment, error)
		Create(context.Context, *Comment) error
	}

	Followers interface {
		Follow(ctx context.Context, followerId, userId int) error
		Unfollow(ctx context.Context, followerId, userId int) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
