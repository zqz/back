package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pressly/chi"
)

// Pagination adds the values "per_page" and "page" to the context with values
// from the query params.
func Pagination(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rawPerPage := chi.URLParam(r, "per_page")
		rawPage := chi.URLParam(r, "page")

		if len(rawPerPage) == 0 {
			rawPerPage = "20"
		}

		if len(rawPage) == 0 {
			rawPage = "0"
		}

		var err error
		var perPage int
		if perPage, err = strconv.Atoi(rawPerPage); err != nil {
			perPage = 20
		}

		var page int
		if page, err = strconv.Atoi(rawPage); err != nil {
			page = 0
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, 1001, perPage)
		ctx = context.WithValue(ctx, 1002, page)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
