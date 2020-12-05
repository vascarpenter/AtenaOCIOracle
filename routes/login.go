package routes

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type loginHTMLtemplate struct {
	Title  string
	NoUser string
	CSS    string
}

// LoginRouter  GET "/login" を処理
func LoginRouter(c echo.Context) error {

	htmlvariable := loginHTMLtemplate{
		Title:  "宛名 管理システム ログイン",
		NoUser: "",
		CSS:    "/css/login.css",
	}
	return c.Render(http.StatusOK, "login", htmlvariable)
}

// LoginRouterPost  POST "/login" を処理
func LoginRouterPost(c echo.Context) error {

	db := Repository()
	defer db.Close()
	ctx := context.Background()

	userID := c.FormValue("userid")
	pass := c.FormValue("password")

	errStr := "指定されたユーザIDが存在しません"

	// SQL: SELECT PASSWORD,ID FROM ATENAUSERS WHERE USERNAME = req.userid limit 1
	var userpass string
	var newuserid int
	err := db.QueryRowContext(ctx, "SELECT PASSWORD,ID FROM ATENAUSERS WHERE USERNAME = :1", userID).Scan(&userpass, &newuserid)

	//	userpasscrypt, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	//	fmt.Printf("ID:%+v,Pass:%+v,Err:%+v,InDBPass:%+v, newCrypt:%+v\n", userID, pass, err, userpass, string(userpasscrypt))

	if err == sql.ErrNoRows {
		//errStr = "指定されたユーザIDが存在しません"
	} else if err != nil {
		panic(err)
	} else if err = bcrypt.CompareHashAndPassword([]byte(userpass), []byte(pass)); err != nil {
		errStr = "パスワードが間違っています"
	} else {
		// login success; create session and redirect to "/"
		session, _ := session.Get("atena_session", c)
		session.Values["userid"] = newuserid
		err = session.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusFound, "/")
	}

	htmlvariable := loginHTMLtemplate{
		Title:  "宛名 管理システム ログイン",
		NoUser: errStr,
		CSS:    "/css/login.css",
	}
	return c.Render(http.StatusOK, "login", htmlvariable)
}
