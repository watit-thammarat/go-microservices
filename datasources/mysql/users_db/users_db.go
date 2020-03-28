package users_db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysql_users_username = "mysql_users_username"
	mysql_users_password = "mysql_users_password"
	mysql_users_host     = "mysql_users_host:3306"
	mysql_users_scheam   = "mysql_users_scheam"
)

var (
	Client *sql.DB
	// username = os.Getenv(mysql_users_username)
	// password = os.Getenv(mysql_users_password)
	// host     = os.Getenv(mysql_users_host)
	// schema   = os.Getenv(mysql_users_scheam)
	username = "root"
	password = "root"
	host     = "localhost:3306"
	schema   = "users_db"
)

func init() {
	datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", username, password, host, schema)
	var err error
	log.Printf("about to connect ot %s", datasourceName)
	Client, err = sql.Open("mysql", datasourceName)
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("database successfully configured")
}
