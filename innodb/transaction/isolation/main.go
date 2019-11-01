package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	createDatabase = `CREATE DATABASE IF NOT EXISTS sample;`
	switchDatabase = `USE sample;`
	createTable    = `CREATE TABLE IF NOT EXISTS simple (
			id    BIGINT UNSIGNED NOT NULL DEFAULT 0,
			value INT UNSIGNED NOT NULL DEFAULT 0,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`
	selectAll = `SELECT * FROM simple;`
)

var (
	db                                                                *sql.DB
	readUncommitTx, readCommittedTx, repeatableReadTx, serializableTx *sql.Tx
)

func init() {
	var err error

	db, err = sql.Open("mysql", "root:single@tcp(127.0.0.1:3306)/mysql")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(createDatabase); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(switchDatabase); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	if readUncommitTx, err = db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	}); err != nil {
		log.Fatal(err)
	}

	if _, err = readUncommitTx.Exec(switchDatabase); err != nil {
		log.Fatal(err)
	}

	if readCommittedTx, err = db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	}); err != nil {
		log.Fatal(err)
	}

	if _, err = readCommittedTx.Exec(switchDatabase); err != nil {
		log.Fatal(err)
	}

	if repeatableReadTx, err = db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	}); err != nil {
		log.Fatal(err)
	}

	if _, err = repeatableReadTx.Exec(switchDatabase); err != nil {
		log.Fatal(err)
	}

	if serializableTx, err = db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}); err != nil {
		log.Fatal(err)
	}

	if _, err = serializableTx.Exec(switchDatabase); err != nil {
		log.Fatal(err)
	}
}

func readAll(label string, tx *sql.Tx) {
	var (
		id    uint64
		value uint32
	)

	log.Printf("--------%s--------", label)
	rows, err := tx.Query(selectAll)
	if err != nil {
		log.Printf("[%s] - Select err %s", label, err)
	}

	for rows.Next() {
		if err = rows.Scan(&id, &value); err != nil {
			log.Fatal(err)
		}
		log.Printf("[%s] id = %d, value = %d", label, id, value)
	}
	log.Println()
}

func writeValue(label string, tx *sql.Tx) {
	log.Printf("--------%s--------\n", label)
	if _, err := tx.Exec(`INSERT INTO simple(id, value) values(1, 100);`); err != nil {
		log.Println(label, err)
	}
	log.Printf("%s Finished\n", label)
}

func readUncommittedWrite() {
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	if err != nil {
		log.Fatal(err)
	}

	if _, err = tx.Exec(switchDatabase); err != nil {
		log.Fatal(err)
	}

	if _, err = tx.Exec(`INSERT INTO simple(id, value) values(1, 100);`); err != nil {
		log.Fatal(err)
	}

	readAll("Read Uncommitted", readUncommitTx)
	readAll("Read Committed", readCommittedTx)
	readAll("Repeatable Read", repeatableReadTx)
	readAll("Serializable", serializableTx)

	if err = tx.Rollback(); err != nil {
		log.Fatal(err)
	}

	log.Println("ok, Rollback() finished")

	readAll("Read Uncommitted", readUncommitTx)
	readAll("Read Committed", readCommittedTx)
	readAll("Repeatable Read", repeatableReadTx)
	readAll("Serializable", serializableTx)
}

func writeAndRollback() {
	wg := &sync.WaitGroup{}
	writeValue("Read Committed", readCommittedTx)

	wg.Add(4)

	go func() {
		time.Sleep(5 * time.Second)
		readCommittedTx.Rollback()
		wg.Done()
	}()

	go func() {
		writeValue("Read Uncommitted", readUncommitTx)
		wg.Done()
	}()

	go func() {
		writeValue("Repeatable Read", repeatableReadTx)
		wg.Done()
	}()

	go func() {
		writeValue("Serializable", serializableTx)
		wg.Done()
	}()

	wg.Wait()
}

func main() {
	writeAndRollback()
}
