package ck
import(
	"github.com/gorilla/securecookie"
	"net/http"
	"log"
	"crypto/sha1"
    "encoding/hex"
	"fmt"
)
//hash generator
func hash(s string) string {
	algorithm := sha1.New()
    algorithm.Write([]byte(s))
    return hex.EncodeToString(algorithm.Sum(nil))
}

//CookieHnadler for cookies
var CookieHandler=securecookie.New(securecookie.GenerateRandomKey(64),securecookie.GenerateRandomKey(32))

//Checkcookie will check for cookies
func Checkcookie(r *http.Request) bool{
	cookie,err := r.Cookie("session") 
	if err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			key:= cookieValue["key"]
			if key!=""{
				return true
			}
		}
	}
	if err!=nil{
		log.Println(err)
	}
	return false
}



//Createcookie will create cookies in the client browser
func Createcookie(w http.ResponseWriter,email string,pass string){
	key:=hash(email+pass)
	value := map[string]string{
		"key": key,
	}
	encoded, err := CookieHandler.Encode("session", value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)		
	}else{
		log.Println(err)
		fmt.Printf("Could not create cookie for email %s",email)
	}
}

//Clearcookie will clear cookies from browser
func Clearcookie(w http.ResponseWriter){
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}