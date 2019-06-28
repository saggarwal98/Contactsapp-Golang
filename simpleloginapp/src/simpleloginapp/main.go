package main
import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	// "time"
	"io/ioutil"
	"github.com/gorilla/securecookie"
	// "dockerpkg/db"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"context"
	"os"
	"os/signal"
)


//CookieHnadler for cookies
var CookieHandler=securecookie.New(securecookie.GenerateRandomKey(64),securecookie.GenerateRandomKey(32))

//main function
func main(){
	handlefunc()
	ctx:=context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()
}

//handlefunc will define all the paths and start the server
func handlefunc(){
	myRouter:=mux.NewRouter()
	myRouter.HandleFunc("/",indexhandler)
	myRouter.HandleFunc("/login",loginhandler)
	myRouter.HandleFunc("/feed",feedhandler)
	myRouter.HandleFunc("/logout",logouthandler)
	myRouter.HandleFunc("/register",registerhandler)
	myRouter.HandleFunc("/addusertodb",addusertodb)
	log.Println("Server started on localhost:3000")
	http.ListenAndServe(":3000",myRouter)
}

//checkcookie will check for cookies
func checkcookie(r *http.Request) bool{
	cookie,err := r.Cookie("session") 
	if err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			email:= cookieValue["email"]
			if email!=""{
				return true
			}
		}
	}
	if err!=nil{
		log.Println(err)
	}
	return false
}



//createcookie will create cookies in the client browser
func createcookie(w http.ResponseWriter,email string){
	value := map[string]string{
		"email": email,
	}
	if encoded, err := CookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}


//clearcookie will clear cookies from browser
func clearcookie(w http.ResponseWriter){
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

//indexhandler will provide space for login
func indexhandler(w http.ResponseWriter,r *http.Request){
	//Check for cookies
	exist:=checkcookie(r)
	if exist==true{
		http.Redirect(w,r,"/feed",301)
	}else{
		content, err := ioutil.ReadFile("simpleloginapp/src/html/login.html")
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintln(w,string(content))
	}
}




//loginhandler will check user details with database and redirect in case of successful authentication
func loginhandler(w http.ResponseWriter,r *http.Request){
	email:=r.FormValue("email")
	password:=r.FormValue("password")
	
	//check in db for user authentication
	str:=checkuser(email,password)
	
	//if user doesn't exist
	if str=="register"{
		content, err := ioutil.ReadFile("simpleloginapp/src/html/register.html")
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintln(w,string(content))
	}else if str=="proceed"{
		//if user exits, save cookie and redirect to feed
		createcookie(w,email)
		http.Redirect(w,r,"/feed",301)
	}else if str=="wrong credentials"{
		http.Redirect(w,r,"/",301)
	}else{
		//some error has been countered
		log.Println(str)
		fmt.Fprintln(w,"Some error has occured. Please try again later")
	}
}


//feedhandler to display user feeds and has a form for user logout
func feedhandler(w http.ResponseWriter,r *http.Request){
	exist:=checkcookie(r)
	if exist==true{
		content, err := ioutil.ReadFile("simpleloginapp/src/html/feed.html")
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintln(w,string(content))
	}else{
		http.Redirect(w,r,"/",301)
	}
}

//logouthandler will clear cookies and redirect to "/"
func logouthandler(w http.ResponseWriter,r *http.Request){
	clearcookie(w)
	content,err:=ioutil.ReadFile("simpleloginapp/src/html/logout.html")
	if err!=nil{
		log.Println(err)
	}else{
		fmt.Fprintf(w,string(content))
	}
}

//registerhandler will be used to register new users
func registerhandler(w http.ResponseWriter,r *http.Request){
	content,err:=ioutil.ReadFile("simpleloginapp/src/html/register.html")
	if err!=nil{
		log.Println(err)
	}
	fmt.Fprintf(w,string(content))
}


//addusertodb adds user to db
func addusertodb(w http.ResponseWriter,r*http.Request){
	// fmt.Fprintln(w,"Your details have been entered into database")
	fname:=r.FormValue("f1")
	log.Println(fname)
	flag:=adduser(r.FormValue("f1"),r.FormValue("l1"),r.FormValue("userEmail"),r.FormValue("userPassword"))
	if flag==false{
		fmt.Println(w,"Internal error")
		fmt.Println(w,"Please try after sometime")
	}else{
		createcookie(w,r.FormValue("userEmail"))
		http.Redirect(w,r,"/feed",301)
	}
}



//**************************************************************
//*******************BELOW SECTION CONTAIN DB FUNCTIONS*********
//**************************************************************

//checkuser function checks for user credentials in database
func checkuser(email,password string) string {
	var db, err= sql.Open("mysql", "saggarwal98:shubham@tcp(0.0.0.0:3306)/simplelogin")
	defer db.Close()
	if err!=nil{
		return "Could not connect to database"
	}
	res1,err:=db.Query("Select * from Users where EMAIL='"+email+"'")
	if err!=nil{
		return "error running query"
	}
	if res1.Next()==true{
		res,err:=db.Query("Select * from Users where EMAIL='"+email+"' AND PASSWORD='"+password+"'")
		if err!=nil{
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


//func adduser will add user to db
func adduser(firstname,lastname,email,password string)bool{
	var db, err= sql.Open("mysql", "saggarwal98:shubham@tcp(0.0.0.0:3306)/simplelogin")
	defer db.Close()
	_,err=db.Query("INSERT INTO Users(FIRSTNAME,LASTNAME,EMAIL,PASSWORD) VALUES('"+firstname+"','"+lastname+"','"+email+"','"+password+"')")
	if err!=nil{
		log.Println(err)
		return false
	}
	return true
}