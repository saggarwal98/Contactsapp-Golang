package db
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
func checkuser(email,password string) string {
	var db, err= sql.Open("mysql", "saggarwal98:shubham@tcp(0.0.0.0:3306)/simplelogin")
	defer db.Close()
	if err!=nil{
		return "Could not connect to database"
	}
	res,err:=db.Query("Select * from Users where EMAIL='"+email+"' AND PASSWORD='"+password+"'")
	if err!=nil{
		return "error running query"
	}else if res.Next()==true{
		return "proceed"
	}else{
		return "register"
	}
}