package routes

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type resetpassHTMLtemplate struct {
	Title       string
	UserName    string
	ErrPassRest string
	CSS         string
}

// ResetPassRouter  GET "/resetpass" を処理
func ResetPassRouter(c echo.Context) error {

	username := c.Get("UserName").(string)
	htmlvariable := resetpassHTMLtemplate{
		Title:       "パスワードの変更",
		UserName:    username,
		ErrPassRest: "",
		CSS:         "/css/resetpass.css",
	}
	return c.Render(http.StatusOK, "resetpass", htmlvariable)
}

// ResetPassRouterPost  POST "/resetpass" を処理
func ResetPassRouterPost(c echo.Context) error {
	userid := c.Get("UserID").(int)
	username := c.Get("UserName").(string)
	htmlvariable := resetpassHTMLtemplate{
		Title:       "パスワードの変更",
		UserName:    username,
		ErrPassRest: "",
		CSS:         "/css/resetpass.css",
	}

	db := Repository()
	defer db.Close()
	ctx := context.Background()

	oldpass := c.FormValue("oldpass")
	newpass := c.FormValue("newpass")

	var userpass string
	err := db.QueryRowContext(ctx, "SELECT PASSWORD FROM ATENAUSERS WHERE ID = :1", userid).Scan(&userpass)

	if err != nil || err == sql.ErrNoRows {
		htmlvariable.ErrPassRest = "ユーザが存在しません"
		return c.Render(http.StatusOK, "resetpass", htmlvariable)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userpass), []byte(oldpass)); err != nil {
		htmlvariable.ErrPassRest = "パスワードが間違っています"
		return c.Render(http.StatusOK, "resetpass", htmlvariable)
	}

	newpasscrypt, _ := bcrypt.GenerateFromPassword([]byte(newpass), 10)

	tx, _ := db.BeginTx(ctx, nil)
	if _, err := tx.Exec("UPDATE ATENAUSERS SET PASSWORD = :1 WHERE ID = :2", string(newpasscrypt), userid); err != nil {
		panic(err)
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	htmlvariable.ErrPassRest = "パスワードが変更されました"
	return c.Render(http.StatusOK, "resetpass", htmlvariable)
}
