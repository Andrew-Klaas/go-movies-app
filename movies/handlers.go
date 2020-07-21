package movies

import (
	"net/http"

	"github.com/GoesToEleven/golang-web-dev/final_project/config"
	"github.com/GoesToEleven/golang-web-dev/final_project/users"
)

//CheckMembership ...
func CheckMembership(w http.ResponseWriter, req *http.Request) {
	u := users.GetUser(w, req)
	if !users.AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if u.Role != "user" {
		http.Error(w, "You must be a user to enter the store", http.StatusForbidden)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
}

//MovieStore ...
func MovieStore(w http.ResponseWriter, req *http.Request) {
	CheckMembership(w, req)

	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	movies, err := AllMovies()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	config.TPL.ExecuteTemplate(w, "moviestore.gohtml", movies)
}

//Show ...
func Show(w http.ResponseWriter, req *http.Request) {
	//CheckMembership(w, req)
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	mv, err := SingleMovie(req)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	config.TPL.ExecuteTemplate(w, "show.gohtml", mv)
}

//Create ...
func Create(w http.ResponseWriter, req *http.Request) {
	//CheckMembership(w, req)
	config.TPL.ExecuteTemplate(w, "create.gohtml", nil)
}

//CreateProcess ...
func CreateProcess(w http.ResponseWriter, req *http.Request) {
	//CheckMembership(w, req)
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}
	mv, err := CreateMovie(req)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	config.TPL.ExecuteTemplate(w, "created.gohtml", mv)
}

//Update ...
func Update(w http.ResponseWriter, req *http.Request) {
	//CheckMembership(w, req)
	mv, err := SingleMovie(req)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	config.TPL.ExecuteTemplate(w, "update.gohtml", mv)
}

//UpdateProcess ...
func UpdateProcess(w http.ResponseWriter, req *http.Request) {
	//CheckMembership(w, req)
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	mv, err := UpdateMovie(req)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	config.TPL.ExecuteTemplate(w, "updated.gohtml", mv)
}

//DeleteProcess ...
func DeleteProcess(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	err := DeleteMovie(req)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	http.Redirect(w, req, "/movies", http.StatusSeeOther)
}
