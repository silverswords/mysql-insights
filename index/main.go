package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/silverswords/mysql-insights/init"
)

func main() {
	hang := make(chan struct{})
	db := init.CreateCon("3306")

	start := time.Now()
	go func() {
		db.CreateIndex()
		log.Println("[addIndex time]", time.Now().Sub(start).Seconds())
	}()
	<-hang
}

func (db *DB) CreateIndex() {
	_, err := db.Exec("ALTER TABLE masterSlaveDB.masterSlaveTable ADD PRIMARY KEY (id);")
	if err != nil {
		log.Println(err)
	}
}
