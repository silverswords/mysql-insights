package main

import (
	"log"
	"strconv"
	"time"

	"github.com/silverswords/mysql-insights/init"
)

func main() {
	hang := make(chan struct{})
	master := init.CreateCon("3306")
	slaveOne := init.CreateCon("3307")
	slaveTwo := init.CreateCon("3308")

	master.CreateDB()
	master.CreateTable()

	slaveOne.CreateDB()
	slaveOne.CreateTable()

	slaveTwo.CreateDB()
	slaveTwo.CreateTable()

	start := time.Now()
	for i := 0; i < 100; i++ {
		master.InsertData(init.Hobby, strconv.Itoa(i))
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
