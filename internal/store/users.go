package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID        int      `json:"id" bson:"id"`
	UserName  string   `json:"username" bson:"username"`
	Email     string   `json:"email" bson:"email"`
	Password  Password `json:"-" bson:"-"`
	CreatedAt string   `json:"created_at" bson:"created_at"`
	IsActive  bool     `json:"is_active" bson:"is_active"`
}

type Password struct {
	Text string
	Hash []byte
}
type UserStore struct {
	db *sql.DB
}

// Note: Remeber one thing once you use a transaction in you code or any function
// then in that function you must have to use transaction not use local database variable or instance.

func (us *UserStore) Create(ctx context.Context, user *User, tx *sql.Tx) error {
	query := `Insert into users (username, password, email)
	values($1, $2, $3) returning id,created_At`

	err := us.db.QueryRowContext(ctx, query,
		user.UserName, user.Password.Hash, user.Email).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (us *UserStore) GetById(ctx context.Context, userId int) (*User, error) {
	query := `
	
	select id, username, email, password, created_at from users where id = $1
	`

	var user User
	err := us.db.QueryRowContext(ctx, query, userId).Scan(&user.ID, &user.UserName, &user.Email, &user.Password.Hash, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (us *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, expiryTime time.Duration) error {

	// Transaction Wrapper.
	// Create the User.
	// Create the User Invite.

	return withTx(us.db, ctx, func(tx *sql.Tx) error {
		if err := us.Create(ctx, user, tx); err != nil {
			return err
		}

		// Create the user invitation.
		if err := us.createUserInvitation(ctx, tx, token, expiryTime, user.ID); err != nil {

			return err
		}

		return nil
	})
}

func (us *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userId int) error {

	query := `insert into users_invitation (token,user_id, expiry) values ($1, $2, $3)`

	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(exp))
	if err != nil {

		fmt.Println("Error - ", err)
		return err
	}

	return nil
}

func (us *UserStore) Activate(ctx context.Context, token string) error {

	// 1. Find the user belongs to the token.
	// 2. Update the User Invitation.
	// 3. Delete the Invitation from the table.

	return withTx(us.db, ctx, func(tx *sql.Tx) error {
		// Step 1:
		user, err := us.GetUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// Step 2:
		user.IsActive = true
		err = us.UpdateUser(ctx, tx, user)
		if err != nil {
			return err
		}

		// Step 3:
		err = us.DeleteUserInvitation(ctx, tx, user.ID)
		if err != nil {
			return err
		}
		return nil
	})

}

func (us *UserStore) GetUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {

	query := `
	
	select  
	u.id, u.username, u.email, u.created_at, u.is_active
	from users as u 
	join 
	users_invitation as ui 
	on ui.user_id = u.id 
	where 
	ui.token = $1 and ui.expiry > $2

	`

	// Convert token to hash string.
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil

}

func (us *UserStore) UpdateUser(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `

        update 
		users set 
		username = $1 , 
		email= $2, 
		is_active = $3 
		where 
		id = $4
    `

	_, err := tx.ExecContext(ctx, query, user.UserName, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) DeleteUserInvitation(ctx context.Context, tx *sql.Tx, userid int) error {
	query := `DELETE FROM users_invitation WHERE user_id = $1`

	_, err := tx.ExecContext(ctx, query, userid)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, userID int) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.DeleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
