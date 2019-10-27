package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

func main() {
	dbmap := initDb("3306")
	defer dbmap.Db.Close()

	// p1 := newPost("aaa", "bbb")
	// start := time.Now()
	// for i := 0; i < 10000; i++ {
	// 	err := dbmap.Insert(&p1)
	// 	log.Println("[current count]", i)
	// 	log.Println("[current insert time]", time.Now().Sub(start).Seconds())
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }

	// start := time.Now()
	// err := dbmap.SelectOne(&p1, "select * from posts where post_id = ?", 10000)
	// log.Println("[query time]", time.Now().Sub(start).Seconds())
	// if err != nil {
	// 	log.Println(err)
	// }

	// update a row
	// start := time.Now()
	// _, err := dbmap.Exec("update posts set title='ccc', body='ddd' where id=9999")
	// log.Println("[update time]", time.Now().Sub(start).Seconds())
	// if err != nil {
	// 	log.Println(err)
	// }

	// err = dbmap.SelectOne(&p1, "select * from posts where post_id=?", p1.Id)
	// checkErr(err, "SelectOne failed")
	// log.Println("p2 row:", p1)

	start := time.Now()
	_, err := dbmap.Exec("delete from posts where id=10000")
	// count, err := dbmap.Delete(&p1)
	log.Println("[delete time]", time.Now().Sub(start).Seconds())
	if err != nil {
		log.Println(err)
	}
	// log.Println("Rows deleted:", count)

	// _, err = dbmap.Exec(c, p1.Id)
	// checkErr(err, "Exec failed")

	// count, err = dbmap.SelectInt("select count(*) from posts")
	// checkErr(err, "select count(*) failed")
	// log.Println("Row count - should be zero:", count)

	// log.Println("Done!")
}

type Post struct {
	Id    int64  `db:"id"`
	Title string `db:"title"`
	Body  string `db:"body"`
}

func newPost(title, body string) Post {
	return Post{
		Title: title,
		Body:  body,
	}
}

func initDb(port string) *gorp.DbMap {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:"+port+")/masterSlaveDB?parseTime=true")
	checkErr(err, "sql.Open failed")

	dialect := gorp.MySQLDialect{"InnoDB", "UTF8"}
	dbmap := &gorp.DbMap{Db: db, Dialect: dialect}

	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
