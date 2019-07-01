package db
import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)
//**************************************************************
//*******************BELOW SECTION CONTAIN DB FUNCTIONS*********
//**************************************************************

//Checkuser function checks for user credentials in database
func Checkuser(email,password string) string {
	var db, err= sql.Open("mysql", "saggarwal98:shubham@tcp(simplelogindatabase:3306)/simplelogin")
	defer db.Close()
	if err!=nil{
		return "Could not connect to database"
	}
	res1,err:=db.Query("Select * from Users where EMAIL='"+email+"'")
	if err!=nil{
		log.Println(err)
		return "error running query"
	}
	if res1.Next()==true{
		res,err:=db.Query("Select * from Users where EMAIL='"+email+"' AND PASSWORD='"+password+"'")
		if err!=nil{
			log.Println(err)
			return "error running query"
		}else if res.Next()==true{
			return "proceed"
		}else{
			return "wrong credentials"
		}
	}else{
		return "register"
	}
}


//func Adduser will add user to db
func Adduser(firstname,lastname,email,password string)bool{
	var db, err= sql.Open("mysql", "saggarwal98:shubham@tcp(simplelogindatabase:3306)/simplelogin")
	log.Println(err)
	defer db.Close()
	_,err=db.Query("INSERT INTO Users(FIRSTNAME,LASTNAME,EMAIL,PASSWORD) VALUES('"+firstname+"','"+lastname+"','"+email+"','"+password+"')")
	if err!=nil{
		log.Println(err)
		return false
	}
	return true
}