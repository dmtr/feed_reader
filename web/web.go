package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dmtr/fbreader/fbutil"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var usertoken = flag.String("usertoken", "", "user token")

func getDB() (sqlx.DB, error) {
	c := getDbConf()
	return sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", c.user, c.dbname))
}

func createSchema(db *sqlx.DB) {
	db.MustExec(schema)
}

var schema = `
CREATE TABLE fbuser (
    id serial,
    username text,
    fbid text
);

CREATE TABLE user_post (
	id text,
	userid FOREIGN KEY id REFERENCES fbuser ON DELETE CASCADE
);

CREATE UNIQUE INDEX user_post_id ON user_post;
`

type dbConf struct {
	user   string
	dbname string
}

func getDbConf() *dbConf {
	c := new(dbConf)
	c.user = os.Getenv("USERNAME")
	c.dbname = os.Getenv("DBNAME")
	return c
}

func main() {
	flag.Parse()
	db, err := getDB()
	if err != nil {
		log.Fatalln(err)
	}

	for r := range fbutil.GetPosts(usertoken) {
		if r.Err != nil {
			log.Printf("got error %s", r.Err)
		} else {
			log.Printf("result %s", r)
		}
	}

}
