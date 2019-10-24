package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	init "github.com/silverswords/mysql-insights/init"
)

func main() {
	hang := make(chan struct{})
	db := init.CreateCon("33")

	start := time.Now()
	go func() {
		db.CreateIndex()
		log.Println("[addIndex time]", time.Now().Sub(start).Seconds())
	}()
	<-hang
}


