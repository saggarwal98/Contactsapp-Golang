package ck
import(
	"net/http"
	// "log"
	// "fmt"
	// "time"
)


//Checkcookie will check for cookies
func Checkcookie(r *http.Request) bool{
	_, err := r.Cookie("session")
	if err==nil{
		return false
	}else{return true}
}



//Createcookie will create cookies in the client browser
func Createcookie(w http.ResponseWriter,email string){
    cookie := http.Cookie{Name: "session", Value: email, MaxAge: 120, Secure:true}
    http.SetCookie(w, &cookie)
}

//Clearcookie will clear cookies from browser
func Clearcookie(w http.ResponseWriter){
	cookie := http.Cookie{Name: "session", Value: "", MaxAge:-1, Secure:true}
    http.SetCookie(w, &cookie)
}



//GetEmailCookie will return email value from cookie
func GetEmailCookie(r *http.Request)string{
	cookie, _ := r.Cookie("session")
	if cookie.Value!=""{
		return cookie.Value
	}else{return "error fetching email"}
}