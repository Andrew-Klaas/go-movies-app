package users

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GoesToEleven/golang-web-dev/final_project/config"
	uuid "github.com/satori/go.uuid"
)

// Index ...
func Index(w http.ResponseWriter, req *http.Request) {
	u := GetUser(w, req)
	//Show()
	config.TPL.ExecuteTemplate(w, "index.gohtml", u)
}

//Signup ...
func Signup(w http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		//grab values from from for user
		un := req.FormValue("username")
		pw := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")
		r := req.FormValue("role")

		if _, ok := Users[un]; ok {
			http.Error(w, "Username Already taken", http.StatusForbidden)
			return
		}

		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = Length
		http.SetCookie(w, c)
		Sessions[c.Value] = Session{un, time.Now()}

		fmt.Printf("Plaintext password: %v\n", pw)
		//HashiCorp Vault encryption
		data := map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(pw)),
		}
		ctxt, err := config.Vclient.Logical().Write("transit/encrypt/my-key", data)
		if err != nil {
			log.Fatal(err)
		}
		s := ctxt.Data["ciphertext"].(string)
		fmt.Printf("Vault encrypted password: %v\n", s)

		Users[un] = User{un, []byte(s), f, l, r}
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	config.TPL.ExecuteTemplate(w, "signup.gohtml", nil)
}

//Login ...
func Login(w http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		pw := req.FormValue("password")

		u, ok := Users[un]
		if !ok {
			http.Error(w, "Username or Password do not match", http.StatusForbidden)
			return
		}

		//HashiCorp Vault decrypt password and check
		data := map[string]interface{}{
			"ciphertext": string(u.Password),
		}
		b64ptxt, err := config.Vclient.Logical().Write("transit/decrypt/my-key", data)
		if err != nil {
			log.Fatal(err)
		}
		s := strings.Split(b64ptxt.Data["plaintext"].(string), ":")
		realptxt, err := base64.StdEncoding.DecodeString(s[0])
		if string(realptxt) != pw {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = Length
		http.SetCookie(w, c)
		Sessions[c.Value] = Session{un, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	config.TPL.ExecuteTemplate(w, "login.gohtml", nil)
}

//Logout ...
func Logout(w http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	c, err := req.Cookie("session")
	if err != nil {
		log.Fatal(err)
	}
	delete(Sessions, c.Value)
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	if time.Now().Sub(LastCleaned) > (time.Second * 30) {
		go Clean()
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

//MovieStore ...
func MovieStore(w http.ResponseWriter, req *http.Request) {
	u := GetUser(w, req)

	if !AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}

	if u.Role != "user" {
		http.Error(w, "You must be a user to enter the store", http.StatusForbidden)
		return
	}
	//Show()
	config.TPL.ExecuteTemplate(w, "moviestore.gohtml", u)
}
