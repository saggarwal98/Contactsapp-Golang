package main
import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"dockerpkg/ck"
	"io/ioutil"
	"dockerpkg/db"
	"context"
	"os"
	"os/signal"
)


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






//indexhandler will provide space for login
func indexhandler(w http.ResponseWriter,r *http.Request){
	//Check for cookies
	exist:=ck.Checkcookie(r)
	if exist==true{
		http.Redirect(w,r,"/feed",301)
	}else{
		content, err := ioutil.ReadFile("./src/html/login.html")
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
	str:=db.Checkuser(email,password)
	
	//if user doesn't exist
	if str=="register"{
		http.Redirect(w,r,"/register",301)
	}else if str=="proceed"{
		//if user exits, save cookie and redirect to feed
		ck.Createcookie(w,email,password)
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
	exist:=ck.Checkcookie(r)
	if exist==true{
		content, err := ioutil.ReadFile("./src/html/feed.html")
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
	ck.Clearcookie(w)
	content,err:=ioutil.ReadFile("./src/html/logout.html")
	if err!=nil{
		log.Println(err)
	}else{
		fmt.Fprintf(w,string(content))
	}
}

//registerhandler will be used to register new users
func registerhandler(w http.ResponseWriter,r *http.Request){
	content,err:=ioutil.ReadFile("./src/html/register.html")
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
	flag:=db.Adduser(r.FormValue("f1"),r.FormValue("l1"),r.FormValue("userEmail"),r.FormValue("userPassword"))
	if flag==false{
		fmt.Println(w,"Internal error")
		fmt.Println(w,"Please try after sometime")
	}else{
		ck.Createcookie(w,r.FormValue("userEmail"),r.FormValue("userPassword"))
		http.Redirect(w,r,"/feed",301)
	}
}