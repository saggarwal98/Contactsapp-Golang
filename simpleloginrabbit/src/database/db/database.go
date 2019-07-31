package db
import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)
//**************************************************************
//*******************BELOW SECTION CONTAIN DB FUNCTIONS*********
//**************************************************************

//func Adduser will add user to db
func Adduser(query string)bool{
	// conn, err := proxysql.NewConn("proxysql", 6033, "proxysql", "proxysqlpassw0rd")
	// if err != nil {
	// 	log.Println(err)
	// }
	// conn.SetCharset("utf8")
	// conn.SetCollation("utf8_general_ci")
	// conn.MakeDBI()
	// db, err := conn.OpenConn()
	db,err:=sql.Open("mysql","saggarwal98:shubham@tcp(localhost:3306)/simplelogin")
	defer db.Close()
	if err!=nil{
		log.Println("Could not connect to database")
	}
	_,err=db.Exec("use simplelogin;")
	_,err=db.Query(query)
	if err!=nil{
		log.Println(err)
		return false
	}
	return true
}
