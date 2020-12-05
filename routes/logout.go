package routes

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
)

// LogoutRouter  handles "/logout"
func LogoutRouter(c echo.Context) error {
	session, _ := session.Get("atena_session", c)
	session.Options.MaxAge = -1 // removes cookie
	session.Values["userid"] = 0
	_ = session.Save(c.Request(), c.Response())
	c.Set("UserID", 0)
	return c.Redirect(http.StatusFound, "/login")
}
