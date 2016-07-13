package users

import "github.com/zqzca/echo"

func Read(c echo.Context) error {
	// tx := StartTransaction()
	// defer tx.Rollback()
	// id := GetParam(c, "id")

	// if u, err := user.FindByID(tx, id); err != nil {
	// 	errors := &UserError{err.Error()}
	// 	return c.JSON(http.StatusOK, u)
	// 	return c.JSON(http.StatusNotFound, errors)
	// } else {
	// 	return c.JSON(http.StatusOK, u)
	// }
	return nil
}
