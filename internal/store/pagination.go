package store

import (
	"fmt"
	"net/http"
	"strconv"
)

// THis is for the pagination data.
type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {

	// THis is for the when the query params are passing like

	// http://localhost:8080/v1/users/feed?limit=2&offset=3
	// now it takes only 2 values out of fetching data,  sorry it fetched only 2 rows. or offset takes 3

	qs := r.URL.Query()

	fmt.Println("QS -", qs)
	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = l
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	return fq, nil
}
