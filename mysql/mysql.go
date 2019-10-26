package mysql

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const Hobby string = "999"

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

func (db *DB) QueryDataByHobbies(hobbies string) (int, error) {
	result := db.QueryRow("SELECT * FROM masterSlaveDB.masterSlaveTable WHERE hobbies = ? LIMIT 1 LOCK IN SHARE MODE", hobbies)
	var (
		name  string
		hobby string
	)
	if err := result.Scan(&name, &hobby); err != nil {
		return 0, err
	}

	return 1, nil
}

func (db *DB) CreateIndex() {
	_, err := db.Exec("CREATE UNIQUE INDEX id ON masterSlaveDB.masterSlaveTable (hobbies);")
	if err != nil {
		log.Println(err)
	}
}

func (db *DB) OrderTable() {
	_, err := db.Exec("SELECT name, hobbies FROM masterSlaveDB.masterSlaveTable ORDER BY hobbies DESC;")
	if err != nil {
		log.Println(err)
	}
}