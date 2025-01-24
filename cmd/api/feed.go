package main

import (
	"fmt"
	"net/http"
	"social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	// Predefined values define for if did not give any query params by user.
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	// If given then parse it and replace the default value.
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	userId := 34

	fmt.Println("Hello ")
	feed, err := app.store.Posts.GetUserFeed(ctx, userId, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJson(w, http.StatusOK, feed)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
