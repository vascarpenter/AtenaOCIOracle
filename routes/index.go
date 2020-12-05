package routes

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

/*
   oracleは''をnullとして扱うので、NullStringやむなし
*/

type tableAtena struct {
	Lastname       string
	Firstname      string
	Furilastname   string
	Furifirstname  string
	Addresscode    string
	Fulladdress    string
	Suffix         string
	Phoneitem      string
	Emailitem      string
	Memo           string
	Namesoffamily1 string
	Suffix1        string
	Namesoffamily2 string
	Suffix2        string
	Namesoffamily3 string
	Suffix3        string
	Atxbaseyear    string
	Nycardhistory  string
	Nycard         string
	Selected       int
	RowNum         string
	ThisYear       string
	LastYear       string
}

// Atenas  slice of tableAtena
type Atenas []tableAtena

type indexHTMLtemplate struct {
	Title    string
	UserName string
	CSS      string
	Iconids  Atenas
}

// getBit make pulldown menu and set default
func getBit(c string, p int, tagname string) string {
	ret1 := "<select name=\"S" + tagname + "\">" // sender = target
	ret2 := "<select name=\"R" + tagname + "\">" // receiver = you
	bitnum, _ := hex.DecodeString("0" + c[p:p+1])
	bitb := bitnum[0]

	if (bitb & 8) != 0 {
		ret1 += `<option value="sou">送</option><option value="mo" selected>喪</option><option value=""></option>`
	} else if (bitb & 4) != 0 {
		ret1 += `<option value="sou" selected>送</option><option value="mo">喪</option><option value=""></option>`
	} else {
		ret1 += `<option value="sou">送</option><option value="mo">喪</option><option value="" selected></option>`
	}
	if (bitb & 2) != 0 {
		ret2 += `<option value="uke">受</option><option value="mo" selected>喪</option><option value=""></option>`
	} else if (bitb & 1) != 0 {
		ret2 += `<option value="uke" selected>受</option><option value="mo">喪</option><option value=""></option>`
	} else {
		ret2 += `<option value="uke">受</option><option value="mo">喪</option><option value="" selected></option>`
	}

	return ret1 + "</select>\n" + ret2 + "</select>\n"
}

// IndexRouter  GET "/" を処理
func IndexRouter(c echo.Context) error {

	db := Repository()
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)

	err := db.PingContext(ctx)
	cancel()
	if err != nil {
		panic(err)
	}

	username := c.Get("UserName").(string)

	orderby := "order by nlssort(FURILASTNAME || FURIFIRSTNAME, 'NLS_SORT=JAPANESE')"
	// cookieがあるか調べる
	if session, err := session.Get("atena_session", c); err == nil {
		if order, ok := session.Values["order"]; ok {
			switch order {
			case "1":
			default:
				orderby = "order by nlssort(FURILASTNAME || FURIFIRSTNAME, 'NLS_SORT=JAPANESE')"
			case "2":
				orderby = "order by selected,nlssort(FURILASTNAME || FURIFIRSTNAME, 'NLS_SORT=JAPANESE')"
			}
		}
	}

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 55*time.Second)
	defer cancel()

	rows, err = db.QueryContext(ctx, `select 
	LASTNAME,  
	FIRSTNAME,  
	FURILASTNAME,  
	FURIFIRSTNAME,  
	ADDRESSCODE,  
	FULLADDRESS,  
	SUFFIX,  
	PHONEITEM,  
	EMAILITEM,  
	MEMO,  
	NAMESOFFAMILY1,  
	SUFFIX1,  
	NAMESOFFAMILY2,  
	SUFFIX2,  
	NAMESOFFAMILY3,  
	SUFFIX3,  
	ATXBASEYEAR,  
	NYCARDHISTORY,  
	NYCARD,
	SELECTED,
	TO_CHAR(ID) from "ADMIN"."ATENA" `+orderby)

	if err != nil {
		panic(err)
	}

	var slice Atenas

	for rows.Next() {
		var oneline tableAtena
		var selectednull sql.NullString
		err = rows.Scan(
			&oneline.Lastname,
			&oneline.Firstname,
			&oneline.Furilastname,
			&oneline.Furifirstname,
			&oneline.Addresscode,
			&oneline.Fulladdress,
			&oneline.Suffix,
			&oneline.Phoneitem,
			&oneline.Emailitem,
			&oneline.Memo,
			&oneline.Namesoffamily1,
			&oneline.Suffix1,
			&oneline.Namesoffamily2,
			&oneline.Suffix2,
			&oneline.Namesoffamily3,
			&oneline.Suffix3,
			&oneline.Atxbaseyear,
			&oneline.Nycardhistory,
			&oneline.Nycard,
			&selectednull,
			&oneline.RowNum,
		)
		if err != nil {
			panic(err)
		}
		oneline.Selected = 0
		if selectednull.Valid && selectednull.String != "" {
			oneline.Selected = 1
		}
		bitstr := oneline.Nycardhistory
		oneline.LastYear = getBit(bitstr, 12, "L"+oneline.RowNum)
		oneline.ThisYear = getBit(bitstr, 13, "T"+oneline.RowNum)

		//		oneline.Nums = int(storage[oneline.ID])
		newslice := append(slice, oneline)
		slice = newslice
		// fmt.Printf("%+v\n", oneline)
	}
	rows.Close()
	//fmt.Printf("%+v %+v\n", slice[1], storage[1])

	htmlvariable := indexHTMLtemplate{
		Title:    "宛名一覧",
		CSS:      "/css/index.css",
		UserName: username,
		Iconids:  slice,
	}

	return c.Render(http.StatusOK, "index", htmlvariable)
}

func calcBits(ststr string, rtstr string) int {
	bit := 0
	if ststr == "mo" {
		bit |= 8
	} else if ststr == "sou" {
		bit |= 4
	}
	if rtstr == "mo" {
		bit |= 2
	} else if rtstr == "uke" {
		bit |= 1
	}
	return bit
}

// IndexRouterPost  POST "/" を処理
func IndexRouterPost(c echo.Context) error {
	db := Repository()
	defer db.Close()
	ctx := context.Background()

	// query max count of atenas
	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ATENA").Scan(&count); err != nil {
		panic(err)
	}

	tx, _ := db.BeginTx(ctx, nil)

	var bit int
	var bit2 int
	for i := 1; i <= count; i++ {
		var oldstr string
		chkstr := fmt.Sprintf("chk%d", i)
		bit = 0
		rtstr := c.FormValue(fmt.Sprintf("RT%d", i))
		ststr := c.FormValue(fmt.Sprintf("ST%d", i))
		rlstr := c.FormValue(fmt.Sprintf("RL%d", i))
		slstr := c.FormValue(fmt.Sprintf("SL%d", i))
		bit = calcBits(ststr, rtstr)
		bit2 = calcBits(slstr, rlstr)

		if err := db.QueryRowContext(ctx, "SELECT NYCARDHISTORY FROM ATENA　WHERE ID = :1", i).Scan(&oldstr); err != nil {
			panic(err)
		}
		newstr := oldstr[:12] + fmt.Sprintf("%x%x", bit2, bit) + oldstr[14:]
		//		fmt.Printf("%d: chk:%s %s %s  '%x%x'\n", i, c.FormValue(chkstr), oldstr, newstr, bit2, bit)

		if _, err := tx.Exec("UPDATE ATENA SET SELECTED=:1, NYCARDHISTORY=:2 WHERE ID =:3", c.FormValue(chkstr), newstr, i); err != nil {
			panic(err)
		}
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return c.Redirect(http.StatusFound, "/")
}
