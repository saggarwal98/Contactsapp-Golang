package main
import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	// "dockerpkg/ck"
	"io/ioutil"
	"dockerpkg/database/db"
	"context"
	"os"
	"os/signal"
	"dockerpkg/rabbit"
	"encoding/base64"
	// "time"
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
	myRouter.HandleFunc("/cContact",createcontacthandler)
	myRouter.HandleFunc("/dContact",deletecontacthandler)
	myRouter.HandleFunc("/eContact",editcontacthandler)
	myRouter.HandleFunc("/getide",idedithandler)
	myRouter.HandleFunc("/completeedit",completeedithandler)
	myRouter.HandleFunc("/completecreate",completecreatehandler)
	myRouter.HandleFunc("/completedelete",completedeletehandler)
	myRouter.HandleFunc("/instantdeletelandingpage",instantdeletelandingpage)
	log.Println("Server started on http://localhost:3000")
	http.ListenAndServe(":3000",myRouter)
}

//Checkcookie will check for cookies
func Checkcookie(r *http.Request) bool{
	_, err := r.Cookie("session")
	if err!=nil{
		return false
	}else{return true}
}



//Createcookie will create cookies in the client browser
func Createcookie(w http.ResponseWriter,email string){
	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))
    cookie := http.Cookie{Name: "session", Value: encodedEmail,MaxAge:86400}
    http.SetCookie(w, &cookie)
}

//Clearcookie will clear cookies from browser
func Clearcookie(w http.ResponseWriter){
	cookie := http.Cookie{Name: "session", Value: "", MaxAge:-1}
    http.SetCookie(w, &cookie)
}



//GetEmailCookie will return email value from cookie
func GetEmailCookie(r *http.Request)string{
	cookie, err := r.Cookie("session")
	if err!=nil{
		return "error fetching email"
	}else{
		email, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err!=nil{
			log.Println(err)
			return "error fetching email"
		}else{
			return string(email)
		}
	}
}

//indexhandler will provide space for login
func indexhandler(w http.ResponseWriter,r *http.Request){
	//Check for cookies
	exist:=Checkcookie(r)
	if exist==true{
		http.Redirect(w,r,"/feed",301)
	}else{
		fmt.Fprintf(w,feed1)
		fmt.Fprintf(w,feed3)
		fmt.Fprintf(w,feed0)
		fmt.Fprintf(w,feed5)
		fmt.Fprintf(w,feed9)
		fmt.Fprintf(w,feed8)
	}
}



//loginhandler will check user details with database and redirect in case of successful authentication
func loginhandler(w http.ResponseWriter,r *http.Request){
	email:=r.FormValue("email")
	password:=r.FormValue("password")
	// log.Println(email)
	
	//check in db for user authentication
	str:=db.Checkuser(email,password)
	
	//if user doesn't exist
	if str=="register"{
		http.Redirect(w,r,"/register",301)
	}else if str=="proceed"{
		//if user exits, save cookie and redirect to feed
		Createcookie(w,email)
		http.Redirect(w,r,"/feed",301)
	}else if str=="wrong credentials"{
		http.Redirect(w,r,"/",301)
	}else{
		//some error has been countered
		//log.Println(str)
		fmt.Fprintln(w,"Some error has occured. Please try again later")
	}
}


//feedhandler to display user feeds and has a form for user logout
func feedhandler(w http.ResponseWriter,r *http.Request){
	exist:=Checkcookie(r)
	if exist==true{
		email:=GetEmailCookie(r)
		if email == "error fetching email"{
			log.Println("could not retrieve email from cookie")
			fmt.Fprintln(w,"Server error")
		}else{
			contacts_id,contacts_name,contacts_phone_number_1,contacts_phone_number_2,contacts_address:=db.GetContacts(email)
			if contacts_id[0]=="No contacts found"{
				fmt.Fprintf(w,feed1)
				fmt.Fprintf(w,feed3)
				fmt.Fprintf(w,feed4)
				fmt.Fprintf(w,feed5)
				fmt.Fprintln(w,`<div>You do not have any contacts</div>`)
				fmt.Fprintf(w,feed8)
			}else{
				fmt.Fprintf(w,feed1)
				fmt.Fprintln(w,feed2)
				fmt.Fprintf(w,feed3)
				fmt.Fprintf(w,feed4)
				fmt.Fprintf(w,feed5)
				fmt.Fprintf(w,feed6)
				for i:=range contacts_id[0:]{
					fmt.Fprintf(w,`<tr id="tb_row"><td id="tb_data1">%s</td><td id="tb_data2">%s</td><td id="tb_data3">%s</td><td id="tb_data4">%s</td><td id="tb_data5">%s</td><td><button onclick="myFunction(this, 'red')">X</button></td></tr>`,contacts_id[i],contacts_name[i],contacts_phone_number_1[i],contacts_phone_number_2[i],contacts_address[i])
				}
				fmt.Fprintf(w,feed7)
				fmt.Fprintf(w,feed8)
			}
		}
	}else{
		http.Redirect(w,r,"/",301)
	}
}

//logouthandler will clear cookies and redirect to "/"
func logouthandler(w http.ResponseWriter,r *http.Request){
	Clearcookie(w)
	fmt.Fprintln(w,`<script>window.onload=()=>{setTimeout(()=>{window.location.href="/"},4000)}</script>`)
	fmt.Fprintln(w,"You have successfully logged out. <br>You will be redirected to login page")
	fmt.Fprintln(w,`You are being redirected to login page. Click here if not redirected <a href="/">Login Page</a>`)
}

//registerhandler will be used to register new users
func registerhandler(w http.ResponseWriter,r *http.Request){
	exist:=Checkcookie(r)
	if exist==true{
		http.Redirect(w,r,"/feed",301)
	}else{
		content,err:=ioutil.ReadFile("./src/html/register.html")
		if err!=nil{
			log.Println(err)
		}
		fmt.Fprintf(w,string(content))
	}
}


//addusertodb adds user to db
func addusertodb(w http.ResponseWriter,r*http.Request){
	// fmt.Fprintln(w,"Your details have been entered into database")
	// fname:=r.FormValue("f1")
	// log.Println(fname)
	exist:=Checkcookie(r)
	if exist==true{
		http.Redirect(w,r,"/feed",301)
	}else{
		query:="INSERT INTO Users(FIRSTNAME,LASTNAME,EMAIL,PASSWORD) VALUES('"+r.FormValue("f1")+"','"+r.FormValue("l1")+"','"+r.FormValue("userEmail")+"','"+r.FormValue("userPassword")+"')"
		rabbit.AddToQueue(query)
		fmt.Fprintln(w,`<script>window.onload=function(){setTimeout(()=>{window.location.href="/"},3000)}</script>`)
		fmt.Fprintf(w,"Your account creation is under process.")
		fmt.Fprintf(w,`You are being redirected to login page. Click here if not redirected <a href="/">Login Page</a>`)
	}
}


//createcontacthandler will discplay form for creatinf contact
func createcontacthandler(w http.ResponseWriter,r *http.Request){
	exist:=Checkcookie(r)
	if exist==true{
		fmt.Fprintf(w,feed1)
		fmt.Fprintf(w,feed3)
		fmt.Fprintf(w,feed4)
		fmt.Fprintf(w,feed5)
		fmt.Fprintf(w,feed12)
		fmt.Fprintf(w,feed8)
	}else{
		http.Redirect(w,r,"/",301)
	}
}



//completecreatehandler will create new contact in database
func completecreatehandler(w http.ResponseWriter,r *http.Request){
	exist:=Checkcookie(r)
	if exist==true{
		email:=GetEmailCookie(r)
		if email == "error fetching email"{
			log.Println("could not retrieve email from cookie")
			fmt.Fprintln(w,"Server error")
		}else{
			var p1,p2 string
			if r.FormValue("phone_number_1")==""{p1="0"}else{p1=r.FormValue("phone_number_1")}
			if r.FormValue("phone_number_2")==""{p2="0"}else{p2=r.FormValue("phone_number_2")}
			rabbit.AddToQueue("INSERT INTO Contacts(Name,Phone_number_1,Phone_number_2,Address,Email) VALUES('"+r.FormValue("name")+"','"+p1+"','"+p2+"','"+r.FormValue("address")+"','"+email+"')")
			http.Redirect(w,r,"/feed",301)
		}
	}else{
		http.Redirect(w,r,"/",301)
	}
}

//editcontacthandler will take id for contact as input
func editcontacthandler(w http.ResponseWriter,r *http.Request){
	exist:=Checkcookie(r)
	if exist==true{
		fmt.Fprintf(w,feed1)
		fmt.Fprintf(w,feed3)
		fmt.Fprintf(w,feed4)
		fmt.Fprintf(w,feed5)
		fmt.Fprintf(w,feed10)
		fmt.Fprintf(w,feed8)
	}else{
		http.Redirect(w,r,"/",301)
	}
}



//func idedithandler will display the edit form
func idedithandler(w http.ResponseWriter,r *http.Request){
	exist:=Checkcookie(r)
	if exist==true{
		email:=GetEmailCookie(r)
		if email == "error fetching email"{
			log.Println("could not retrieve email from cookie")
			fmt.Fprintln(w,"Server error")
		}else{
			enteredid:=r.FormValue("getidfield")
			id,n1,p1,p2,a1:=db.GetContactForEdit(enteredid,email)
			if id!=enteredid{
				fmt.Fprintln(w,`<html><head><script>window.onload=()=>{setTimeout(()=>{window.location.href="/feed"},3000)}</script></head><body>No Contact found with that id.<br>Redirecting you to Contacts page.<a href="/feed">Click here</a> if not redirected</body></html>`)
			}else{
				fmt.Fprintf(w,feed1)
				fmt.Fprintf(w,feed3)
				fmt.Fprintf(w,feed4)
				fmt.Fprintf(w,feed5)
				fmt.Fprintf(w,feed11,id,n1,p1,p2,a1)
				fmt.Fprintf(w,feed8)
			}
		}
	}else{
		http.Redirect(w,r,"/",301)
	}
}



//completeedithandler will update in the database
func completeedithandler(w http.ResponseWriter,r *http.Request){
	exist:=Checkcookie(r)
	if exist==true{
		email:=GetEmailCookie(r)
		if email == "error fetching email"{
			log.Println("could not retrieve email from cookie")
			fmt.Fprintln(w,"Server error")
		}else{
			var p1,p2 string
			if r.FormValue("phone_number_1")==""{p1="0"}else{p1=r.FormValue("phone_number_1")}
			if r.FormValue("phone_number_2")==""{p2="0"}else{p2=r.FormValue("phone_number_2")}
			rabbit.AddToQueue("UPDATE Contacts SET Name='"+r.FormValue("name")+"',Phone_number_1="+p1+",Phone_number_2="+p2+",Address='"+r.FormValue("address")+"' WHERE ID="+r.FormValue("id")+" AND Email='"+email+"'")
			http.Redirect(w,r,"/feed",301)
		}
	}else{
		http.Redirect(w,r,"/",301)
	}
}



//deletecontacthandler will cdisplay form to delete contact
func deletecontacthandler(w http.ResponseWriter,r *http.Request){
	exist:=Checkcookie(r)
	if exist==true{
		fmt.Fprintf(w,feed1)
		fmt.Fprintf(w,feed3)
		fmt.Fprintf(w,feed4)
		fmt.Fprintf(w,feed5)
		fmt.Fprintf(w,feed13)
		fmt.Fprintf(w,feed8)
	}else{
		http.Redirect(w,r,"/",301)
	}
}



//completedeletehandler will delete contact
func completedeletehandler(w http.ResponseWriter,r *http.Request){
	exist:=Checkcookie(r)
	if exist==true{
		email:=GetEmailCookie(r)
		if email == "error fetching email"{
			log.Println("could not retrieve email from cookie")
			fmt.Fprintln(w,"Server error")
		}else{
			count:=db.GetContactForDelete(r.FormValue("id"),r.FormValue("name"),email)
			if count==1{
				fmt.Fprintf(w,`<html><head><script>window.onload=()=>{setTimeout(()=>{window.location.href="/feed"},3000)}</script></head><body>Contact Deleted.<br>Redirecting you to Contacts page.<a href="/feed">Click here</a> if not redirected</body></html>`)
				rabbit.AddToQueue("DELETE FROM Contacts WHERE ID="+r.FormValue("id")+" AND Name='"+r.FormValue("name")+"' AND Email='"+email+"'")
			}else{
				fmt.Fprintf(w,`<html><head><script>window.onload=()=>{setTimeout(()=>{window.location.href="/feed"},3000)}</script></head><body>No Contact found with those details.<br>Redirecting you to Contacts page.<a href="/feed">Click here</a> if not redirected</body></html>`)
			}
		}
	}else{
		http.Redirect(w,r,"/",301)
	}
}


//instantdeletelandingpage will redirect user to contacts page after few seconds
func instantdeletelandingpage(w http.ResponseWriter,r * http.Request){
	fmt.Fprintf(w,`<html><head><script>window.onload=()=>{setTimeout(()=>{window.location.href="/feed"},1500)}</script></head><body>Contact Deleted.<br>Redirecting you to Contacts page.<a href="/feed">Click here</a> if not redirected</body></html>`)
}




var feed1 string=`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
	<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
`
var feed2 string=`<script>
function myFunction(elmnt,clr) {
  elmnt.style.color = clr;
  Id=elmnt.parentElement.parentElement.children[0].innerHTML
  Name=elmnt.parentElement.parentElement.children[1].innerHTML
  $.post('/completedelete',{id:Id,name:Name},function(){window.alert("Contact Deleted");setTimeout(function(){window.location.href="/feed"},200)})
}
</script>`
var feed3 string=`</head>
<body>
	<div class="container-fluid">
        <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
            <div class="collapse navbar-collapse" id="navbarNav">
			    <ul class="navbar-nav">
			 	    <li class="nav-item active">
					 	<a class="nav-link" href="/feed">Display Contacts</a>
					</li>
                    <li class="nav-item active">
                        <a class="nav-link" href="/cContact">Create Contact</a>
                    </li>
                    <li class="nav-item active">
                        <a class="nav-link" href="/eContact">Edit Contact</a>
                    </li>
                    <li class="nav-item active">
                        <a class="nav-link" href="/dContact">Delete Contact</a>
                    </li>        
                </ul>
			</div>`
var feed4 string=`<form class="form-inline" action="/logout" method="POST">
    			<button class="btn btn-outline-success my-2 my-sm-0" type="submit">Logout</button>
			  </form>`
var feed0 string=`<form class="form-inline" action="/register" method="GET">
			  <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Register</button>
			</form>`
var feed5 string=`</nav>`
var feed6 string=`<h1>Contacts</h1>
		  <hr style="display: block;border-width: 4px;background:grey;height:10px;">
            <table class="table" id="table1">
                    <thead>
                      <tr>
                        <th scope="col">ID</th>
                        <th scope="col">Name</th>
                        <th scope="col">Phone Number 1</th>
                        <th scope="col">Phone Number 2</th>
                        <th scope="col">Address</th>
                      </tr>
                    </thead>
                    <tbody>`
var feed7 string=`</tbody></table><hr>`     
var feed8 string=`</div></body></html>`
var feed9 string=`<br><br><h3>Please Login to continue</h3><br><form  action="/login" method="POST"><div class="form-group"><label>Email address</label>
<input type="email" class="form-control" id="email" name="email" placeholder="Email"></div><div class="form-group">
<label>Password</label><input type="password" class="form-control" id="password" name="password" placeholder="Password">
</div><button type="submit" class="btn btn-primary">Login</button></form>`
var feed10 string=`<form action="/getide" method="POST"><div class="form-group"><label>Enter the ID of contact</label>
<input type="number" class="form-control" id="getidfield" name="getidfield" placeholder="Enter contact id">
</div><button type="submit">SUBMIT</button></form>`
var feed11 string=`<form  action="/completeedit" method="POST">
<div class="form-group"><label>ID</label>
<input type="text" class="form-control" id="id" name="id" value="%s"></div><div class="form-group">
<label>Name</label><input type="text" class="form-control" id="name" name="name" value="%s">
</div><div class="form-group"><label>Phone_number_1</label>
<input type="number" class="form-control" id="phone_number_1" name="phone_number_1" value="%s"></div>
<div class="form-group"><label>Phone_number_2</label>
<input type="number" class="form-control" id="phone_number_2" name="phone_number_2" value="%s"></div>
<div class="form-group"><label>Address</label>
<input type="text" class="form-control" id="address" name="address" value="%s"></div>
<button type="submit" class="btn btn-primary">Update Details</button></form><br><button onclick="window.location.href='/feed'">Cancel</button><br> <b>Please do not change ID else your data may be lost</b>`
var feed12 string=`<form  action="/completecreate" method="POST"><div class="form-group">
<label>Name</label><input type="text" class="form-control" id="name" name="name">
</div><div class="form-group"><label>Phone_number_1</label>
<input type="number" class="form-control" id="phone_number_1" name="phone_number_1"></div>
<div class="form-group"><label>Phone_number_2</label>
<input type="number" class="form-control" id="phone_number_2" name="phone_number_2"></div>
<div class="form-group"><label>Address</label>
<input type="text" class="form-control" id="address" name="address"></div>
<button type="submit" class="btn btn-primary">Create Contact</button></form><br><button onclick="window.location.href='/feed'">Cancel</button>`
var feed13 string=`<form action="/completedelete" method="POST"><div class="form-group"><label>Enter the ID of contact</label>
<input type="number" class="form-control" id="id" name="id" placeholder="Enter contact id">
</div><div class="form-group"><label>Enter the Name of contact</label>
<input type="text" class="form-control" id="name" name="name" placeholder="Enter contact name">
</div><button type="submit">Confirm Delete</button></form><br><button onclick="window.location.href='/feed'">Cancel</button>`