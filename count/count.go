package main

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const hobby string = "99"

func main() {
	hang := make(chan struct{})
	master := CreateCon("3306")
	slaveOne := CreateCon("3307")
	slaveTwo := CreateCon("3308")

	master.CreateDB()
	master.CreateTable()

	slaveOne.CreateDB()
	slaveOne.CreateTable()

	slaveTwo.CreateDB()
	slaveTwo.CreateTable()

	start := time.Now()
	for i := 0; i < 100; i++ {
		master.InsertData(hobby, strconv.Itoa(i))
		log.Println("[current count]", i)
		log.Println("[current insert time]", time.Now().Sub(start).Seconds())
	}

	// start = time.Now()
	// go func() {
	// 	var code int
	// 	for code != 1 {
	// 		code = master.QueryDataByHobbies(hobby)
	// 	}
	// 	log.Println("[master query]", time.Now().Sub(start).Seconds())
	// }()

	// start = time.Now()
	// go func() {
	// 	var code int
	// 	for code != 1 {
	// 		code = slaveOne.QueryDataByHobbies(hobby)
	// 	}
	// 	log.Println("[slave 3307 query]", time.Now().Sub(start).Seconds())
	// }()

	// start = time.Now()
	// go func() {
	// 	var code int
	// 	for code != 1 {
	// 		code = slaveTwo.QueryDataByHobbies(hobby)
	// 	}
	// 	log.Println("[slave 3308 query]", time.Now().Sub(start).Seconds())
	// }()

	<-hang
}

type DB struct {
	*sql.DB
}

// CreateCon create a db conn to local with given port
func CreateCon(port string) *DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:"+port+")/?parseTime=true")
	if err != nil {
		log.Println(err)
	}

	return &DB{
		db,
	}
}

func (db *DB) CreateDB() error {
	result, err := db.Exec("CREATE DATABASE IF NOT EXISTS masterSlaveDB;")
	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.New("[create order] : create order affected 0 rows")
	}
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (db *DB) CreateTable() {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS masterSlaveDB.masterSlaveTable (" +
			// "id bigint(20) unsigned NOT NULL AUTO_INCREMENT, " +
			"name varchar(50) DEFAULT NULL, " +
			"hobbies varchar(200) DEFAULT NULL " +
			// "PRIMARY KEY (id)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; ")
	if err != nil {
		log.Println(err)
	}
}

func (db *DB) InsertData(name, hobbies string) {
	_, err := db.Exec("INSERT INTO masterSlaveDB.masterSlaveTable (name, hobbies) VALUES (?,?)", name, hobbies)
	if err != nil {
		log.Println(err)
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
