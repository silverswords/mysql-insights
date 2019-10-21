package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const hobby string = "999"

func main() {
	hang := make(chan struct{})
	master := CreateCon("3306")
	slaveOne := CreateCon("3307")
	slaveTwo := CreateCon("3308")

	start := time.Now()
	go func() {
		var code int
		for code != 1 {
			code = master.QueryDataByHobbies(hobby)
		}
		log.Println("[master query]", time.Now().Sub(start).Seconds())
	}()

	start = time.Now()
	go func() {
		var code int
		for code != 1 {
			code = slaveOne.QueryDataByHobbies(hobby)
		}
		log.Println("[slave 3307 query]", time.Now().Sub(start).Seconds())
	}()

	start = time.Now()
	go func() {
		var code int
		for code != 1 {
			code = slaveTwo.QueryDataByHobbies(hobby)
		}
		log.Println("[slave 3308 query]", time.Now().Sub(start).Seconds())
	}()

	<-hang
}

type DB struct {
	*sql.DB
}

func CreateCon(port string) *DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:"+port+")/?parseTime=true")
	if err != nil {
		log.Println(err)
	}

	return &DB{
		db,
	}
}

func (db *DB) QueryDataByHobbies(hobbies string) int {
	result := db.QueryRow("SELECT * FROM masterSlaveDB.masterSlaveTable WHERE name = ? LIMIT 1 LOCK IN SHARE MODE", hobbies)
	var (
		id    int64
		name  string
		hobby string
	)
	if err := result.Scan(&id, &name, &hobby); err != nil {
		return 0
	}

	return 1
}
