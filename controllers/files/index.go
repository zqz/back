package files

import "net/http"

//Index returns a list of files
func (f Controller) Index(w http.ResponseWriter, r *http.Request) {
	// pageStr := e.QueryParam("page")
	// perPageStr := e.QueryParam("per_page")

	// page, err := strconv.Atoi(pageStr)
	// if err != nil {
	// 	f.Warn("Failed to decode page param", "e", err)
	// 	return e.NoContent(http.StatusBadRequest)
	// }
	// perPage, err := strconv.Atoi(perPageStr)
	// if err != nil {
	// 	f.Warn("Failed to decode per_page param")
	// 	return e.NoContent(http.StatusBadRequest)
	// }

	// files, err := models.Files(f.DB, qm.OrderBy("created_at desc"), qm.Limit(perPage), qm.Offset(page*perPage)).All()
	// if err != nil {
	// 	fmt.Println("failed to fetch page:", err)
	// 	return err
	// }

	// return e.JSON(http.StatusOK, files)
}
