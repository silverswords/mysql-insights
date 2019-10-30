package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

func main() {
	hang := make(chan struct{})
	dbmap := initDb("3307")
	// defer dbmap.Db.Close()

	// insert
	// p1 := newPost("aaa", "bbb")
	// start := time.Now()
	// for i := 0; i < 10; i++ {
	// 	go func() {
	// 		hang := make(chan struct{})
	// 		log.Println("sss")
	// 		err := dbmap.Insert(&p1)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		<-hang
	// 	}()
	// 	// err := dbmap.Insert(&p1)
	// 	log.Println("[current count]", i)
	// 	// log.Println("[current insert time]", time.Now().Sub(start).Seconds())
	// }

	// // query
	// start := time.Now()
	// err := dbmap.SelectOne(&p1, "select * from posts where id = ?", 100000)
	// log.Println("[query time]", time.Now().Sub(start).Seconds())
	// if err != nil {
	// 	log.Println(err)
	// }

	// // update a row
	// start = time.Now()
	// _, err = dbmap.Exec("update posts set title='ccc', body='ddd' where id=100000")
	// log.Println("[update time]", time.Now().Sub(start).Seconds())
	// if err != nil {
	// 	log.Println(err)
	// }

	// // update all row
	// start = time.Now()
	// _, err = dbmap.Exec("update posts set title='ccc', body='ddd'")
	// log.Println("[update time]", time.Now().Sub(start).Seconds())
	// if err != nil {
	// 	log.Println(err)
	// }

	// // delete all row
	// start = time.Now()
	// _, err = dbmap.Exec("delete from posts")
	// log.Println("[delete time]", time.Now().Sub(start).Seconds())
	// if err != nil {
	// 	log.Println(err)
	// }

	per := &Person{0, 0, 0, "bob", "smith"}
	inv := &Invoice{0, 0, 0, "xmas order", per.Id}
	go InsertInv(dbmap, inv, per)
	<-hang
}

// type Post struct {
// 	Id    int64  `db:"id"`
// 	Title string `db:"title"`
// 	Body  string `db:"body"`
// }

// func newPost(title, body string) Post {
// 	return Post{
// 		Title: title,
// 		Body:  body,
// 	}
// }

func initDb(port string) *gorp.DbMap {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:"+port+")/masterSlaveDB?parseTime=true")
	checkErr(err, "sql.Open failed")

	dialect := gorp.MySQLDialect{"InnoDB", "UTF8"}
	dbmap := &gorp.DbMap{Db: db, Dialect: dialect}

	dbmap.AddTableWithName(Person{}, "person").SetKeys(true, "Id")
	dbmap.AddTableWithName(Invoice{}, "invoice").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

// transactions
func InsertInv(dbmap *gorp.DbMap, inv *Invoice, per *Person) error {
	trans, err := dbmap.Begin()
	checkErr(err, "Begin failed")

	err = trans.Insert(per)
	checkErr(err, "Insert failed")

	inv.PersonId = per.Id
	err = trans.Insert(inv)
	checkErr(err, "Insert failed")

	return trans.Commit()
}

type Person struct {
	Id      int64
	Created int64
	Updated int64
	FName   string
	LName   string
}

type Invoice struct {
	Id       int64
	Created  int64
	Updated  int64
	Memo     string
	PersonId int64
}
