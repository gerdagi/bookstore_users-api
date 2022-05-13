package users_db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb" // driver for mssql
)

const (
	mssql_users_username = "kforms"
	mssql_users_password = "JA}XRr5*9mYhCWcXwwjg!"
	mssql_users_host     = "20.224.77.124"
	mssql_users_schema   = "AlternatifForms"
)

var (
	Client *sql.DB

	username = mssql_users_username //os.Getenv(mssql_users_username)
	password = mssql_users_password //os.Getenv(mssql_users_password)
	host     = mssql_users_host     //os.Getenv(mssql_users_host)
	schema   = mssql_users_schema   //os.Getenv(mssql_users_schema)
)

func init() {

	datasourceName := fmt.Sprintf("database=%s;server=%s;password=%s;user id=%s;", schema,
		host,
		password,
		username)

	var err error
	Client, err = sql.Open("mssql", datasourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	log.Println("database successfully configured")
}
