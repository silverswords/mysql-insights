package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/silverswords/mysql-insights/mysql"
)

func main() {
	hang := make(chan struct{})
	db := sql.CreateCon("3306")

	start := time.Now()
	go func() {
		db.OrderTable()
		log.Println("[Order time]", time.Now().Sub(start).Seconds())
	}()
	<-hang
}
