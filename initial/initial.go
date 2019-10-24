package init

import (
	"database/sql"
	"errors"
	"log"

)

const Hobby string = "99"

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

func (db *DB) CreateIndex() {
	_, err := db.Exec("ALTER TABLE masterSlaveDB.masterSlaveTable ADD PRIMARY KEY (id);")
	if err != nil {
		log.Println(err)
	}
}
