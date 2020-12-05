package routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type registerHTMLtemplate struct {
	Title       string
	EmailExists string
	CSS         string
}

// RegisterRouter  GET "/register" を処理
func RegisterRouter(c echo.Context) error {

	htmlvariable := registerHTMLtemplate{
		Title:       "ユーザー登録",
		EmailExists: "",
		CSS:         "/css/register.css",
	}
	return c.Render(http.StatusOK, "register", htmlvariable)
}

// RegisterRouterPost  POST "/register" を処理
func RegisterRouterPost(c echo.Context) error {
	db := Repository()
	defer db.Close()
	ctx := context.Background()

	htmlvariable := registerHTMLtemplate{
		Title:       "ユーザー登録",
		EmailExists: "",
		CSS:         "/css/register.css",
	}

	userName := c.FormValue("username")
	userpass := c.FormValue("userpass")
	mailaddress := c.FormValue("mailaddress")

	// 'SELECT * FROM ATENAUSERS WHERE USERNAME = ? LIMIT 1'
	var usercnt int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ATENAUSERS WHERE USERNAME = :1", userName).Scan(&usercnt)

	if err == nil && usercnt > 0 {
		htmlvariable.EmailExists = "同じユーザーIDがすでに存在します"
		return c.Render(http.StatusOK, "register", htmlvariable)
	}

	userpasscrypt, err := bcrypt.GenerateFromPassword([]byte(userpass), 10)
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ATENAUSERS").Scan(&usercnt); err != nil {
		panic(err)
	}

	var storage []byte
	storage = make([]byte, 300)
	for i := 0; i < 300; i++ {
		storage[i] = 0
	}
	tx, _ := db.BeginTx(ctx, nil)
	if _, err := tx.Exec(`INSERT INTO ATENAUSERS (ID,USERNAME,PASSWORD,STORAGE,BINDING,NUMCHILD)
		values (:1,:2,:3,:4,:5,:6)`,
		usercnt+1,             // 1: ID
		userName,              // 2: USERNAME
		string(userpasscrypt), // 3: userpass
		"",                    // 4: STORAGE: initial null string
		mailaddress,           // 5: binding
		storage,               // 6: numchild
	); err != nil {
		panic(err)
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
	htmlvariable.EmailExists = "追加しました"
	return c.Render(http.StatusOK, "register", htmlvariable)
}
