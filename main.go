package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type test_board struct {
	Bno     int
	Name    string
	Age     int
	Content string
	Regdate string
}

type PassedData struct {
	PostData []test_board
}

// DB connection info
const (
	host     = "DESKTOP-JDH1ELA.local"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "postgres"
)

// DB connection
func getDBConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode = disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	fmt.Println("Successfully created connection to database")
	return db
}

// error
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////		/////////		////////////////////////	/////////////////	////////	/////////////	//////////////////////////////
// ///////////	////	/////////			///////////////////	////	////////////////////////	/	/////////	//////////////////////////////
// /////////	/////	///////// //////	////////////////	////////	/////////	/////////	////	//////	//////////////////////////////
// ////////		////////	/// ////////	///////////////	////////////	/////////	////////	///////	/////	//////////////////////////////
// //////	  /////////////	/ /////////////	////////////	///////////////	/////////	/////////	////////	//	//////////////////////////////
// ////		///////////////	///////////////	///////////						/////////	////////	////////		//////////////////////////////
// ////		////////////// ////////////////	/////////	/////////////////	/////////	////////	////////////	//////////////////////////////
// /////	////////////// ////////////////	////		////////////////////	////	////////	////////////	//////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// main
func main() {
	e := echo.New()
	e.Static("static", "web/static")

	//load template from folder
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("web/template/*.html")),
	}

	// get, post????????? ?????? ????????????
	e.Renderer = renderer
	e.GET("/", getHome)                            // ????????????
	e.GET("/boardList", getBoardList)              //  ????????? ??????
	e.GET("/boardregistration", getRegistration)   //????????? ???????????????
	e.POST("/boardregistration", postRegistration) // ????????? ??????
	e.GET("/boardDetail", getDetail)               // ???????????????
	e.GET("/boardUpdate", getUpdate)               // ???????????????
	e.POST("/boardUpdate", postUpdate)             // ????????????
	e.POST("/boardDetail", postDelete)             // ????????????

	e.Logger.Fatal(e.Start(":1111"))

}

// ?????????
func getHome(c echo.Context) error {
	// fmt.Println("getHome??????")
	return c.Render(http.StatusOK, "home.html", "")
}

// ????????? ??????
func getBoardList(c echo.Context) error {

	return c.Render(http.StatusOK, "boardList.html", getTempData())
}

// ????????? ????????? ????????????
func getTempData() *PassedData {

	// fmt.Println("getTempData() ??????")

	// db??????
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)

	datas := []test_board{}
	sql_statement := "select * from test_board order by bno ASC" // ????????? ??????
	rows, err := db.Query(sql_statement)                         // ?????? ????????? ??????
	CheckError(err)                                              // ?????? ??????
	defer rows.Close()
	for rows.Next() { // for????????? row??? ????????? ??????
		data := test_board{}
		rows.Scan(&data.Bno, &data.Name, &data.Age, &data.Content, &data.Regdate)
		datas = append(datas, data)
	}

	// fmt.Println(datas)
	temp := PassedData{
		PostData: datas,
	}
	return &temp
}

// ????????? ???????????????
func getRegistration(c echo.Context) error {
	return c.Render(http.StatusOK, "boardregistration.html", "")
}

// ????????? ????????????
func postRegistration(c echo.Context) error {

	// db??????
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)
	// fmt.Println("postRegistration() ??????")

	// form submit??? ?????? ???????????? ???????????? ??????
	name := c.Request().FormValue("Name")
	age := c.Request().FormValue("Age")
	content := c.Request().FormValue("Content")

	// ????????? ??????Query???
	sql_statementInsert := "insert into "
	sql_statementInsert += "test_board(name, age, content, regdate) "
	sql_statementInsert += "values ($1, $2, $3, now());"
	_, err = db.Exec(sql_statementInsert, name, age, content)
	CheckError(err)
	fmt.Println("Insert 1 row of data")
	return c.Render(http.StatusOK, "boardList.html", getTempData())
}

// ???????????????
func getDetail(c echo.Context) error {
	Bno := c.Request().FormValue("Bno")

	// db??????
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)

	// bno select Query???
	sql_statementInsert := "select * from test_board where bno = $1"
	rows, err := db.Query(sql_statementInsert, Bno)
	CheckError(err)
	defer rows.Close() //????????? ????????? (???????????? ??????)

	data := test_board{}
	for rows.Next() {
		err := rows.Scan(&data.Bno, &data.Name, &data.Age, &data.Content, &data.Regdate)
		CheckError(err)
	}

	// fmt.Println(data)
	return c.Render(http.StatusOK, "boardDetail.html", data)
}

// ???????????????
func getUpdate(c echo.Context) error {
	Bno := c.Request().FormValue("Bno")

	// db??????
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)

	// bno select Query???
	sql_statementInsert := "select * from test_board where bno = $1"
	rows, err := db.Query(sql_statementInsert, Bno)
	CheckError(err)
	defer rows.Close() //????????? ????????? (???????????? ??????)

	data := test_board{}
	for rows.Next() {
		err := rows.Scan(&data.Bno, &data.Name, &data.Age, &data.Content, &data.Regdate)
		CheckError(err)
	}
	return c.Render(http.StatusOK, "boardUpdate.html", data)
}

func postUpdate(c echo.Context) error {

	// db??????
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)
	// fmt.Println("postRegistration() ??????")

	// form submit??? ?????? ???????????? ???????????? ??????
	bno := c.Request().FormValue("Bno")
	name := c.Request().FormValue("Name")
	age := c.Request().FormValue("Age")
	content := c.Request().FormValue("Content")

	// fmt.Println("name" + name + "age" + age + "content" + content + "bno" + bno)

	// ????????? ?????? Query???
	sql_statementUpdate := "update test_board "
	sql_statementUpdate += "set name = $1, age = $2, content = $3 "
	sql_statementUpdate += "where bno = $4"
	_, err = db.Exec(sql_statementUpdate, name, age, content, bno)
	CheckError(err)
	fmt.Println("Update 1 row of data")
	return c.Render(http.StatusOK, "boardList.html", getTempData())
}

// ????????? ??????
func postDelete(c echo.Context) error {
	// db??????
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)
	// fmt.Println("postDelete() ??????")

	// form submit??? ?????? ???????????? ???????????? ??????
	bno := c.Request().FormValue("Bno")

	// fmt.Println("bno ==> " + bno)

	// ????????? ?????? Query???
	sql_statementDelete := "delete from test_board "
	sql_statementDelete += "where bno = $1"
	_, err = db.Exec(sql_statementDelete, bno)
	CheckError(err)
	fmt.Println("Delete 1 row of data")
	return c.Render(http.StatusOK, "boardList.html", getTempData())
}
