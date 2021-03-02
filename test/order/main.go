package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/silverswords/mysql-insights/mysql"
)

func main() {
	const hobby string = "999999hobbies"
	db := sql.CreateCon("3306")

	start := time.Now()
	db.OrderTable()
	log.Println("[OrderQuery time]", time.Now().Sub(start).Seconds())
}
