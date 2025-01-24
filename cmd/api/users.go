package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

var userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	Id := chi.URLParam(r, "userId")
	userId, err := strconv.Atoi(Id)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := app.store.Users.GetById(ctx, userId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)

		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = writeJson(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type FollowUser struct {
	UserId int `json:"userId"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)
	// Id := chi.URLParam(r, "userId")
	// followedId, err := strconv.Atoi(Id)
	// if err != nil {
	// 	app.badRequestError(w, r, err)
	// 	return
	// }

	var payload FollowUser
	if err := readJosn(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followerUser.ID, payload.UserId); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	UnfollowerUser := getUserFromContext(r)
	// Id := chi.URLParam(r, "userId")
	// followedId, err := strconv.Atoi(Id)
	// if err != nil {
	// 	app.badRequestError(w, r, err)
	// 	return
	// }

	var payload FollowUser
	if err := readJosn(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, UnfollowerUser.ID, payload.UserId); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
			return
		case store.ErrNotFound:
			app.internalServerError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			Id := chi.URLParam(r, "userId")
			userId, err := strconv.Atoi(Id)
			if err != nil {
				app.badRequestError(w, r, err)
				return
			}

			ctx := r.Context()
			user, err := app.store.Users.GetById(ctx, userId)
			if err != nil {
				switch {
				case errors.Is(err, store.ErrNotFound):
					app.notFoundResponse(w, r, err)

				default:
					app.internalServerError(w, r, err)
				}
				return
			}
			ctx = context.WithValue(ctx, userCtx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = writeJson(w, http.StatusOK, "")
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
