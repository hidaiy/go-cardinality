package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

type User struct {
	Id        int
	name      string
	email     string
	sex       int
	createdAt time.Time
}

const timestampFormat = "2006-01-02 15:04:05"

type Config struct {
	user     string
	password string
	host     string
	port     int
	database string
}

func craeteDBConnectString(c Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.user, c.password, c.host, c.port, "mysql")
}

func main() {
	config := Config{
		user:     "root",
		password: "root",
		host:     "127.0.0.1",
		port:     3306,
		database: "mydatabase",
	}

	// DB接続
	db, err := sql.Open("mysql", craeteDBConnectString(config))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	defer db.Close()

	rows, err := db.Query("select * from users")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	defer rows.Close()

	var user User
	var createdAt string

	for rows.Next() {
		rows.Scan(&user.Id, &user.name, &user.email, &user.sex, &createdAt)
		createdAtTime, err := time.Parse(timestampFormat, createdAt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		user.createdAt = createdAtTime
		fmt.Printf("%#v time: %s\n", user, user.createdAt.Format(timestampFormat))
	}
}
