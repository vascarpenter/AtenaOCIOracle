package main

import (
	"AtenaOCIOracle/m/routes"
	"context"
	"database/sql"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	_ "github.com/mattn/go-oci8"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// to create models: "sqlboiler --wipe mysql"

// Template はHTMLテンプレートを利用するためのRenderer Interfaceです。
type Template struct {
	templates *template.Template
}

// Render はHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込みます。
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// redirectLoginWithoutAuth  contextにIDが入っていないか 0であった場合は、login画面にリダイレクトする
func redirectLoginWithoutAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userid := c.Get("UserID")
			if userid == 0 || userid == nil {
				// not login'd, go to login page
				return c.Redirect(http.StatusFound, "/login")
			}
			return next(c)
		}
	}
}

// setUserMiddleware　cookieを参照して、ユーザがログインしていればdbにアクセスし名前、IDをcontextに入れる
func setUserMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			db := routes.Repository()
			defer db.Close()

			ctx := context.Background()

			if session, err := session.Get("atena_session", c); err == nil {
				if userid, ok := session.Values["userid"]; ok {
					// cookie exist
					if useridint := userid.(int); useridint > 0 {
						// get UserName name
						var username string
						err := db.QueryRowContext(ctx, "SELECT USERNAME FROM ATENAUSERS WHERE ID = :1", useridint).Scan(&username)
						//fmt.Printf("%+v\n", hospname)

						switch {
						case err == sql.ErrNoRows:
							// cannot get; pass through
						case err != nil:
							panic(err)
						default:
							c.Set("UserName", username)
							c.Set("UserID", useridint)
						}
					}
				}
			}

			return next(c)
		}
	}
}

func main() {
	// Echo instance
	e := echo.New()

	funcMap := template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
		"minus":    func(a, b int) int { return a - b },
		"mod":      func(a, b int) int { return a % b },
	}
	t := &Template{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} uri=${uri} path=${path} status=${status}\n",
	}))
	e.Use(middleware.Recover())

	var keystore string
	if keystore = os.Getenv("COOKIE_SEED"); keystore == "" {
		keystore = "secret randseed key"
	}
	var store = sessions.NewCookieStore([]byte(keystore))
	e.Use(session.Middleware(store))

	// Routes
	e.Static("/css", "./static/css")
	e.Static("/img", "./static/img")
	e.Static("/javascript", "./static/javascript")
	e.Static("/icon", "./static/icon")
	e.GET("/", routes.IndexRouter, setUserMiddleware(), redirectLoginWithoutAuth())
	e.POST("/", routes.IndexRouterPost, setUserMiddleware(), redirectLoginWithoutAuth())
	e.GET("/login", routes.LoginRouter)
	e.POST("/login", routes.LoginRouterPost)
	e.GET("/logout", routes.LogoutRouter)
	//e.GET("/register", routes.RegisterRouter) // 現在１ユーザのみ；登録不可
	//e.POST("/register", routes.RegisterRouterPost)
	e.GET("/resetpass", routes.ResetPassRouter, setUserMiddleware(), redirectLoginWithoutAuth())
	e.POST("/resetpass", routes.ResetPassRouterPost, setUserMiddleware(), redirectLoginWithoutAuth())
	e.GET("/listby/:id", routes.ListByRouter, setUserMiddleware(), redirectLoginWithoutAuth())

	// Start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3002"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
