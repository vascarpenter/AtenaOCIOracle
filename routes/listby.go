package routes

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
)

type numsSlices struct {
	ID   int
	Nums int
}

// ListByRouter  handles "/listby/:id"
func ListByRouter(c echo.Context) error {
	id := c.Param("id")

	session, _ := session.Get("atena_session", c)
	session.Values["order"] = id
	_ = session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/")
}
