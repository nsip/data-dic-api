package entity

import (
	"net/http"

	lk "github.com/digisan/logkit"
	"github.com/labstack/echo/v4"
)

// @Title insert or update one entity data
// @Summary insert or update one entity data by a json file
// @Description
// @Tags    Ingest
// @Accept  application/octet-stream
// @Produce json
// @Param   entity path string true "entity name for incoming entity data"
// @Success 200 "OK - insert or update successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/upsert [post]
func UpsertEntity(c echo.Context) error {

	lk.Log("...")

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

	return nil

}
