package users

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

//Length ...
const Length int = 3000

//Sessions SessionID|Session{UserName,LastActivity}
var Sessions = map[string]Session{}

//Users UserName|User
var Users = map[string]User{}

var LastCleaned time.Time = time.Now()

//GetUser ...
func GetUser(w http.ResponseWriter, req *http.Request) User {
	c, err := req.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
	}
	c.MaxAge = Length
	http.SetCookie(w, c)

	var u User
	if s, ok := Sessions[c.Value]; ok {
		s.LastActivity = time.Now()
		Sessions[c.Value] = s
		u = Users[s.UserName]
	}
	return u
}

//AlreadyLoggedIn checks if a sessions exists, if it does not then return false.
//If it does contain a session return if the user exists: true/false. Also update the cookie
func AlreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	s, ok := Sessions[c.Value]
	if ok {
		s.LastActivity = time.Now()
		Sessions[c.Value] = s
	}
	_, ok = Users[s.UserName]
	c.MaxAge = Length
	http.SetCookie(w, c)
	return ok
}

//Show ...
func Show() {
	fmt.Println("******SESSIONS MAP******\n")
	for k, v := range Sessions {
		fmt.Printf("SessionID: %v, Session: %v\n", k, v)
	}
	fmt.Println("********\n")
}

func Clean() {
	for k, v := range Sessions {
		if time.Now().Sub(v.LastActivity) > (time.Second * 3000) {
			delete(Sessions, k)
		}
	}
	Show()
}
