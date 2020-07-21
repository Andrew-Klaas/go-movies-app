package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
	_ "github.com/lib/pq"
)

//DB Connection
var DB *sql.DB
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

//VClient ...
var Vclient, _ = api.NewClient(&api.Config{Address: "http://127.0.0.1:8200", HttpClient: httpClient})
var tokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
var k8sAuthRole = "go-movies-app"
var k8sAuthPath = "auth/kubernetes/login"

func init() {

	//K8s auth
	buf, err1 := ioutil.ReadFile(tokenPath)
	if err1 != nil {
		log.Fatal(err1)
	}
	jwt := string(buf)
	fmt.Printf("k8s Service Account jwt: %v\n", jwt)
	options := map[string]interface{}{
		"jwt":  jwt,
		"role": k8sAuthRole,
	}
	secret, err1 := Vclient.Logical().Write(k8sAuthPath, options)
	if err1 != nil {
		log.Fatal(err1)
	}
	token := secret.Auth.ClientToken
	fmt.Printf("Vault token: %v\n", token)
	Vclient.SetToken(token)

	//Vault dynamic DB secrets
	data, _ := Vclient.Logical().Read("database/creds/my-role")
	username := data.Data["username"]
	password := data.Data["password"]
	SQLQuery := "postgres://" + username.(string) + ":" + password.(string) + "@localhost/movies?sslmode=disable"

	//Open DB connection
	var err error
	fmt.Printf("SQLQUERY: %v\n", SQLQuery)
	DB, err = sql.Open("postgres", SQLQuery)
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")
	SQLQuery = "CREATE TABLE movies ( movieID SERIAL PRIMARY KEY,title TEXT NOT NULL, director TEXT NOT NULL, price REAL DEFAULT 25500.00 );"
	_, err = DB.Exec(SQLQuery)
	SQLQuery = "INSERT INTO movies (movieID, title,director,price) VALUES(1,'Gladiator', 'Ridley Scott', 10.99);"
	_, err = DB.Exec(SQLQuery)
	SQLQuery = "INSERT INTO movies (movieID, title,director,price) VALUES(2,'Kingdom of Heaven', 'Ridley Scott', 10.99);"
	_, err = DB.Exec(SQLQuery)
	SQLQuery = "INSERT INTO movies (movieID, title,director,price) VALUES(3,'Fast and the Furious', 'Rob Cohen', 10.99);"
	_, err = DB.Exec(SQLQuery)
}
