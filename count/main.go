package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/silverswords/mysql-insights/mysql"
)

func main() {
	hang := make(chan struct{})
	master := sql.CreateCon("3306")
	slaveOne := sql.CreateCon("3307")
	slaveTwo := sql.CreateCon("3308")

	master.CreateDB()
	master.CreateTable()

	slaveOne.CreateDB()
	slaveOne.CreateTable()

	slaveTwo.CreateDB()
	slaveTwo.CreateTable()

	// start := time.Now()
	// for i := 0; i < 1000; i++ {
	// 	master.InsertData(sql.Hobby, strconv.Itoa(i))
	// 	log.Println("[current count]", i)
	// 	log.Println("[current insert time]", time.Now().Sub(start).Seconds())
	// }
	start := time.Now()
	go func() {
		var code int
		var err error
		for code != 1 {
			code, err = master.QueryDataByHobbies(sql.Hobby)
			if err != nil {
				log.Println(err)
			}
		}
		log.Println("[master query]", time.Now().Sub(start).Seconds())
	}()

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
