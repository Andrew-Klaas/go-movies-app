package config

/*this will likely need to be moved in order to use Vault*/

import (
	"database/sql"
	"fmt"
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
var k8sAuthRole = "demo"
var k8sAuthPath = "auth/kubernetes/login"

func init() {

	/*
		buf, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			log.Fatal(err)
		}
		jwt := string(buf)

		options := map[string]interface{}{
			"jwt":  jwt,
			"role": k8sAuthRole,
		}
		secret, err := Vclient.Logical().Write(k8sAuthPath, options)
		if err != nil {
			log.Fatal(err)
		}

		token := secret.Auth.ClientToken
	*/
	token := "password"
	Vclient.SetToken(token)

	data, _ := Vclient.Logical().Read("database/creds/my-role")
	fmt.Printf("data %v\n", data)
	username := data.Data["username"]
	password := data.Data["password"]
	fmt.Printf("username %v\n password %v\n", username, password)
	SQLQuery := "postgres://" + username.(string) + ":" + password.(string) + "@localhost/movies?sslmode=disable"

	//Open DB connection
	var err error
	fmt.Printf("SQLQUERY: %v\n", SQLQuery)
	DB, err = sql.Open("postgres", SQLQuery)
	//DB, err = sql.Open("postgres", "postgres://v-token-my-role-9bElU3kj9pyyrUyQMLwt-1595268254:A1a-wMngz6Iza5AmWmFH@localhost/movies?sslmode=disable")
	//goapp / password
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")
	/*
		SQLQuery := "CREATE TABLE movies ( movieID SERIAL PRIMARY KEY,title TEXT NOT NULL, director TEXT NOT NULL, price REAL DEFAULT 25500.00 );"
		_, err = DB.Exec(SQLQuery)
		SQLQuery = "INSERT INTO movies (movieID, title,director,price) VALUES(1,'Gladiator', 'Ridley Scott', 10.99);"
		_, err = DB.Exec(SQLQuery)
	*/
}

/*
	#transit setup
	vault secrets enable transit
	vault write -f transit/keys/my-key

	##database setup
	vault secrets enable database

	vault write database/config/my-postgresql-database \
    plugin_name=postgresql-database-plugin \
    allowed_roles="my-role" \
    connection_url="postgresql://{{username}}:{{password}}@localhost:5432/movies?sslmode=disable" \
    username="goapp" \
	password="password"

	vault write database/roles/my-role \
    db_name=my-postgresql-database \
    creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
    ALTER USER \"{{name}}\" WITH SUPERUSER;" \
    default_ttl="1h" \
    max_ttl="24h"

	vault read database/creds/my-role

*/
