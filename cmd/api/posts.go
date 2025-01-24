package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   string `json:"title" validate:"omitempty,max=100"`
	Content string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var postPayload CreatePostPayload

	err := readJosn(w, r, &postPayload)
	if err != nil {
		fmt.Println("errro getting the read json data ", err)
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(postPayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	userId := 1

	post := &store.Post{
		Title:   postPayload.Title,
		Content: postPayload.Content,
		Tags:    postPayload.Tags,

		// Todo change after auth.
		UserID: int(userId),
	}
	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	postId, err := strconv.Atoi(idParam)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	post, err := app.store.Posts.GetById(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.internalServerError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	fmt.Println("Post ------ ")

	comments, err := app.store.Comments.GetByPostId(ctx, postId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletepostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	postId, err := strconv.Atoi(idParam)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	err = app.store.Posts.Delete(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)

		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postId")
		postId, err := strconv.Atoi(idParam)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		post, err := app.store.Posts.GetById(ctx, postId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.internalServerError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getFromContext(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Inside the update")
	post := getFromContext(r)

	fmt.Println("Update post ", post)

	var updatePayload UpdatePostPayload

	if err := readJosn(w, r, &updatePayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	fmt.Println("Read Post - ", updatePayload)

	if err := Validate.Struct(updatePayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if updatePayload.Title == "" || updatePayload.Title != "" {
		post.Title = updatePayload.Title
	}

	if updatePayload.Content == "" || updatePayload.Content != "" {
		post.Content = updatePayload.Content
	}

	ctx := r.Context()
	if err := app.store.Posts.Update(ctx, post); err != nil {

		fmt.Println("Error -", err)
		return
	}

	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
