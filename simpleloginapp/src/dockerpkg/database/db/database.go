package db
import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	// "github.com/imSQL/proxysql"
	// "encoding/json"
	// "reflect"
)
//**************************************************************
//*******************BELOW SECTION CONTAIN DB FUNCTIONS*********
//**************************************************************

type Contact struct {
	ID             string `json:"ID"`
	Name           string `json:"Title"`
	Phone_number_1 int    `json:"Phone_number_1"`
	Phone_number_2 int    `json:"Phone_number_2"`
	Address        string `json:"Address"`
	Email          string `json:"Email"`
}


//Checkuser function checks for user credentials in database
func Checkuser(email,password string) string {
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
		return "Could not connect to database"
	}
	res1,err:=db.Query("Select * from Users where EMAIL='"+email+"'")
	if err!=nil{
		log.Println(err)
		log.Println("abcdef")
		return "error running query"
	}
	if res1.Next()==true{
		res,err:=db.Query("Select * from Users where EMAIL='"+email+"' AND PASSWORD='"+password+"'")
		if err!=nil{
			log.Println(err)
			log.Println("abcdefg")
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
	_,err=db.Exec(query)
	if err!=nil{
		log.Println(err)
		log.Println("abc1")
		return false
	}
	return true
}



//func GetContacts will return contacts of each user
func GetContacts(email string)([]string,[]string,[]string,[]string,[]string) {
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
	results, err := db.Query("Select * from Contacts WHERE Email='"+email+"'")
	flag := false
	var p1,p2 string
	var id,name,phone_number_1,phone_number_2,address []string
	if results.Next() == true {
		flag = true
		var a Contact
		err = results.Scan(&a.ID, &a.Name, &a.Phone_number_1, &a.Phone_number_2,&a.Address,&a.Email)
		if err != nil {
			log.Println(err)
			log.Println("abc2")
			return id,name,phone_number_1,phone_number_2,address
		}
		if a.Phone_number_1==0{p1=""}else{p1=strconv.Itoa(a.Phone_number_1)}
		if a.Phone_number_2==0{p2=""}else{p2=strconv.Itoa(a.Phone_number_2)}
		id = append(id,a.ID)
		name=append(name,a.Name)
		phone_number_1=append(phone_number_1,p1)
		phone_number_2=append(phone_number_2,p2)
		address=append(address,a.Address)
		for results.Next() {
			err = results.Scan(&a.ID, &a.Name, &a.Phone_number_1, &a.Phone_number_2,&a.Address,&a.Email)
			if err != nil {
				log.Println(err)
				log.Println("abc3")
				return id,name,phone_number_1,phone_number_2,address
			}
			if a.Phone_number_1==0{p1=""}else{p1=strconv.Itoa(a.Phone_number_1)}
			if a.Phone_number_2==0{p2=""}else{p2=strconv.Itoa(a.Phone_number_2)}
			id = append(id,a.ID)
			name=append(name,a.Name)
			phone_number_1=append(phone_number_1,p1)
			phone_number_2=append(phone_number_2,p2)
			address=append(address,a.Address)
		}
	}
	if flag == false {
		id=append(id,"No contacts found")
	}
	return id,name,phone_number_1,phone_number_2,address
}


//GetContactForEdit
func GetContactForEdit(id,email string)(string,string,string,string,string){
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
	results,err:=db.Query("SELECT ID,Name,Phone_number_1,Phone_number_2,Address FROM Contacts WHERE ID="+ id +" AND Email='"+ email+"'")
	if err!=nil{
		log.Println(err)
		log.Println("abc4")
	}
	if results.Next() == true {
		var a Contact
		err = results.Scan(&a.ID,&a.Name,&a.Phone_number_1,&a.Phone_number_2,&a.Address)
		if err != nil {
			log.Println(err)
			log.Println("abc5")
		}
	return a.ID,a.Name,strconv.Itoa(a.Phone_number_1),strconv.Itoa(a.Phone_number_2),a.Address
	}
	return "0",err.Error(),"","",""
}




//GetContactForDelete will return presence of contact in database
func GetContactForDelete(id,name,email string)int{
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
	results,err:=db.Query("SELECT * FROM Contacts WHERE ID="+ id +" AND Email='"+ email+"' AND Name='"+name+"'")
	if err!=nil{
		log.Println(err)
	}
	if results.Next() == true {
	return 1
	}
	return 0
}
