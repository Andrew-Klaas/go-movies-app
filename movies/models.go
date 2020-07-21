package movies

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/GoesToEleven/golang-web-dev/final_project/config"
)

// Movie ...
type Movie struct {
	MovieID  string
	Title    string
	Director string
	Price    float32
}

//AllMovies ...
func AllMovies() ([]Movie, error) {
	var err error
	//var rows *sql.Rows
	rows, err := config.DB.Query("SELECT * FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mvs := make([]Movie, 0)
	for rows.Next() {
		mv := Movie{}
		err := rows.Scan(&mv.MovieID, &mv.Title, &mv.Director, &mv.Price)
		if err != nil {
			return nil, err
		}
		mvs = append(mvs, mv)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return mvs, nil
}

//SingleMovie ...
func SingleMovie(req *http.Request) (Movie, error) {
	mv := Movie{}
	mvID := req.FormValue("movieID")
	if mvID == "" {
		return mv, errors.New("400. Bad Request")
	}
	row := config.DB.QueryRow("SELECT * FROM movies WHERE movieID = $1", mvID)
	err := row.Scan(&mv.MovieID, &mv.Title, &mv.Director, &mv.Price)
	if err != nil {
		return mv, err
	}

	return mv, nil
}

//CreateMovie ...
func CreateMovie(req *http.Request) (Movie, error) {
	mv := Movie{}
	mv.MovieID = req.FormValue("movieID")
	mv.Title = req.FormValue("title")
	mv.Director = req.FormValue("director")
	price := req.FormValue("price")

	if mv.MovieID == "" || mv.Title == "" || mv.Director == "" || price == "" {
		return mv, errors.New("400. Bad request. All fields must be complete")
	}

	// convert form values
	f64, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return mv, errors.New("406. Not Acceptable. Price must be a number")
	}
	mv.Price = float32(f64)

	_, err = config.DB.Exec("INSERT INTO movies (movieID,title,director,price) VALUES ($1, $2, $3, $4)", mv.MovieID, mv.Title, mv.Director, mv.Price)
	if err != nil {
		return mv, errors.New("500. Internal Server Error." + err.Error())
	}

	return mv, nil
}

//UpdateMovie ...
func UpdateMovie(req *http.Request) (Movie, error) {
	mv := Movie{}
	mv.MovieID = req.FormValue("movieID")
	mv.Title = req.FormValue("title")
	mv.Director = req.FormValue("director")
	price := req.FormValue("price")

	if mv.MovieID == "" || mv.Title == "" || mv.Director == "" || price == "" {
		return mv, errors.New("400. Bad request. All fields must be complete")
	}

	// convert form values
	f64, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return mv, errors.New("406. Not Acceptable. Price must be a number")
	}
	mv.Price = float32(f64)

	_, err = config.DB.Exec("UPDATE movies SET movieID = $1, title=$2, director=$3, price=$4 WHERE movieID=$1;", mv.MovieID, mv.Title, mv.Director, mv.Price)
	if err != nil {
		return mv, err
	}
	return mv, nil
}

//DeleteMovie ...
func DeleteMovie(req *http.Request) error {
	mvID := req.FormValue("movieID")

	_, err := config.DB.Exec("DELETE FROM movies WHERE movieID=$1;", mvID)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil

}
