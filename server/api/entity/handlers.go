package entity

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nsip/data-dic-api/server/api/db"
)

// @Title insert or update one entity data
// @Summary insert or update one entity data by a json file
// @Description
// @Tags    Ingest
// @Accept  application/octet-stream
// @Produce json
// @Param   entity path string true "entity name for incoming entity data"
// @Success 200 "OK - insert or update successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/insert [post]
func Insert(c echo.Context) error {

	var (
		entity  = c.Param("entity")
		dataRdr = c.Request().Body
	)

	if len(entity) == 0 {
		return c.String(http.StatusBadRequest, "entity name is empty")
	}
	if dataRdr == nil {
		return c.String(http.StatusBadRequest, "body data is empty")
	}

	db.UseDbCol("testing", "entity")
	id, err := db.Insert(dataRdr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, id)
}

// @Router /api/entity/find [get]
func Find(c echo.Context) error {
	return nil
}
