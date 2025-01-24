package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"social/internal/mailer"
	"social/internal/store"
	"social/internal/utils"

	"github.com/google/uuid"
)

type UserPayload struct {
	UserName string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserWithTOken struct {
	store.User
	Token string `json:"token" bson:"token"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload UserPayload

	err := readJosn(w, r, &payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	passwrdHash, err := utils.GenerateHash(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user := &store.User{
		UserName: payload.UserName,
		Email:    payload.Email,
		Password: passwrdHash,
	}

	// Token will be as byte or hash.
	// 1. Make a hash of that text.
	// 2. Convert into the string of that hash.

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	ctx := r.Context()

	err = app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithTOken := UserWithTOken{
		User:  *user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.UserName,
		ActivationURL: activationURL,
	}

	// send mail
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.UserName, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent", "status code", status)

	if err := writeJson(w, http.StatusCreated, userWithTOken); err != nil {
		app.internalServerError(w, r, err)
	}

	fmt.Println("Successfully store in the DB")

	err = writeJson(w, http.StatusCreated, userWithTOken)
	if err != nil {
		app.internalServerError(w, r, err)
	}

}
