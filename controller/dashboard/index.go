package dashboard

import (
	"database/sql"
	"log"
	"math"
	"net/http"
	"time"

	"gopkg.in/nullbio/null.v4"

	"github.com/pressly/chi/render"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/dependencies"
)

type dashboardEntry struct {
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	ThumbnailID null.String `json:"thumb_id"`
	CreatedAt   time.Time   `json:"created_at"`
}

// Controller is a exposed struct
type Controller struct {
	dependencies.Dependencies
}

type dashboardData struct {
	Entries *[]dashboardEntry `json:"data"`
	Page    int               `json:"current_page"`
	Total   int               `json:"total_pages"`
}

const paginationSQL = `
	SELECT
	f.name, t.id, f.slug, f.created_at
	FROM files AS f
	LEFT JOIN thumbnails as t
	ON t.file_id = f.id
	ORDER BY f.created_at DESC
	OFFSET $1
	LIMIT $2
`

const totalPagesSQL = `SELECT count(*) FROM files`

// Index returns a list of files
func (c Controller) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	perPage := ctx.Value(1001).(int)
	page := ctx.Value(1002).(int)

	entries, err := pagination(c.DB, page, perPage)

	if err != nil {
		c.Error("failed to fetch page:", "err", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	total := totalPages(c.DB, perPage)

	data := dashboardData{
		Entries: entries,
		Total:   total,
		Page:    page,
	}

	render.JSON(w, r, data)
}

func totalPages(ex db.Executor, perPage int) int {
	var count int

	err := ex.QueryRow(totalPagesSQL).Scan(&count)

	if err != nil {
		return 0
	}

	return int(math.Ceil(float64(count) / float64(perPage)))
}

func pagination(ex db.Executor, page int, perPage int) (*[]dashboardEntry, error) {
	var entries []dashboardEntry
	var err error
	var rows *sql.Rows

	offset := perPage * page

	if rows, err = ex.Query(paginationSQL, offset, perPage); err != nil {
		return &entries, err
	}
	defer rows.Close()

	for rows.Next() {
		var e dashboardEntry

		err = rows.Scan(
			&e.Name, &e.ThumbnailID, &e.Slug, &e.CreatedAt,
		)

		if err != nil {
			log.Fatal(err)
		}

		entries = append(entries, e)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return &entries, err
}
