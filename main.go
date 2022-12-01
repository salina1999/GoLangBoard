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

	// get, post요청이 오면 함수실행
	e.Renderer = renderer
	e.GET("/", getHome)                            // 실행화면
	e.GET("/boardList", getBoardList)              //  게시물 목록
	e.GET("/boardregistration", getRegistration)   //게시물 등록페이지
	e.POST("/boardregistration", postRegistration) // 게시물 등록
	e.GET("/boardDetail", getDetail)               // 상세페이지
	e.GET("/boardUpdate", getUpdate)               // 수정페이지
	e.POST("/boardUpdate", postUpdate)             // 수정하기
	e.POST("/boardDetail", postDelete)             // 삭제하기

	e.Logger.Fatal(e.Start(":1111"))

}

// 홈화면
func getHome(c echo.Context) error {
	// fmt.Println("getHome시작")
	return c.Render(http.StatusOK, "home.html", "")
}

// 게시물 목록
func getBoardList(c echo.Context) error {

	return c.Render(http.StatusOK, "boardList.html", getTempData())
}

// 게시판 리스트 불러오기
func getTempData() *PassedData {

	// fmt.Println("getTempData() 시작")

	// db연결
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)

	datas := []test_board{}
	sql_statement := "select * from test_board order by bno ASC" // 쿼리문 작성
	rows, err := db.Query(sql_statement)                         // 전달 인자에 넣음
	CheckError(err)                                              // 에러 체크
	defer rows.Close()
	for rows.Next() { // for문으로 row행 레코드 읽음
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

// 게시물 등록페이지
func getRegistration(c echo.Context) error {
	return c.Render(http.StatusOK, "boardregistration.html", "")
}

// 게시물 등록하기
func postRegistration(c echo.Context) error {

	// db연결
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)
	// fmt.Println("postRegistration() 시작")

	// form submit을 통해 전달되는 파라미터 획득
	name := c.Request().FormValue("Name")
	age := c.Request().FormValue("Age")
	content := c.Request().FormValue("Content")

	// 게시물 등록Query문
	sql_statementInsert := "insert into "
	sql_statementInsert += "test_board(name, age, content, regdate) "
	sql_statementInsert += "values ($1, $2, $3, now());"
	_, err = db.Exec(sql_statementInsert, name, age, content)
	CheckError(err)
	fmt.Println("Insert 1 row of data")
	return c.Render(http.StatusOK, "boardList.html", getTempData())
}

// 상세페이지
func getDetail(c echo.Context) error {
	Bno := c.Request().FormValue("Bno")

	// db연결
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)

	// bno select Query문
	sql_statementInsert := "select * from test_board where bno = $1"
	rows, err := db.Query(sql_statementInsert, Bno)
	CheckError(err)
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	data := test_board{}
	for rows.Next() {
		err := rows.Scan(&data.Bno, &data.Name, &data.Age, &data.Content, &data.Regdate)
		CheckError(err)
	}

	// fmt.Println(data)
	return c.Render(http.StatusOK, "boardDetail.html", data)
}

// 수정페이지
func getUpdate(c echo.Context) error {
	Bno := c.Request().FormValue("Bno")

	// db연결
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)

	// bno select Query문
	sql_statementInsert := "select * from test_board where bno = $1"
	rows, err := db.Query(sql_statementInsert, Bno)
	CheckError(err)
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	data := test_board{}
	for rows.Next() {
		err := rows.Scan(&data.Bno, &data.Name, &data.Age, &data.Content, &data.Regdate)
		CheckError(err)
	}
	return c.Render(http.StatusOK, "boardUpdate.html", data)
}

func postUpdate(c echo.Context) error {

	// db연결
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)
	// fmt.Println("postRegistration() 시작")

	// form submit을 통해 전달되는 파라미터 획득
	bno := c.Request().FormValue("Bno")
	name := c.Request().FormValue("Name")
	age := c.Request().FormValue("Age")
	content := c.Request().FormValue("Content")

	// fmt.Println("name" + name + "age" + age + "content" + content + "bno" + bno)

	// 게시물 수정 Query문
	sql_statementUpdate := "update test_board "
	sql_statementUpdate += "set name = $1, age = $2, content = $3 "
	sql_statementUpdate += "where bno = $4"
	_, err = db.Exec(sql_statementUpdate, name, age, content, bno)
	CheckError(err)
	fmt.Println("Update 1 row of data")
	return c.Render(http.StatusOK, "boardList.html", getTempData())
}

// 게시물 삭제
func postDelete(c echo.Context) error {
	// db연결
	db := getDBConnection()
	defer db.Close()
	err := db.Ping()
	CheckError(err)
	// fmt.Println("postDelete() 시작")

	// form submit을 통해 전달되는 파라미터 획득
	bno := c.Request().FormValue("Bno")

	// fmt.Println("bno ==> " + bno)

	// 게시물 수정 Query문
	sql_statementDelete := "delete from test_board "
	sql_statementDelete += "where bno = $1"
	_, err = db.Exec(sql_statementDelete, bno)
	CheckError(err)
	fmt.Println("Delete 1 row of data")
	return c.Render(http.StatusOK, "boardList.html", getTempData())
}
