package main

import (
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/silverswords/mysql-insights/mysql"
)

func main() {
	hang := make(chan struct{})
	const name string = "singing"
	// const hobby string = "99singing"
	const count = 100000
	master := sql.CreateCon("3306")
	slaveOne := sql.CreateCon("3307")
	slaveTwo := sql.CreateCon("3308")

	master.CreateDB()
	master.CreateTable()

	slaveOne.CreateDB()
	slaveOne.CreateTable()

	slaveTwo.CreateDB()
	slaveTwo.CreateTable()

	start := time.Now()
	for i := 0; i < count; i++ {
		master.InsertData(name, strconv.Itoa(i)+"singing")
		log.Println("[current count]", i)
		log.Println("[current insert time]", time.Now().Sub(start).Seconds())
	}

	go func() {
		var code int
		var err error
		for code != 1 {
			code, err = master.QueryDataByHobbies()
			if err != nil {
				log.Println(err)
			}
		}
		log.Println("[master 3306 query]", time.Now().Sub(start).Seconds())
	}()

	go func() {
		var code int
		var err error
		for code != 1 {
			code, err = slaveOne.QueryDataByHobbies()
			if err != nil {
				log.Println(err)
			}
		}
		log.Println("[slave 3307 query]", time.Now().Sub(start).Seconds())
	}()

	go func() {
		var code int
		var err error
		for code != 1 {
			code, err = slaveTwo.QueryDataByHobbies()
			if err != nil {
				log.Println(err)
			}
		}
		log.Println("[slave 3308 query]", time.Now().Sub(start).Seconds())
	}()

	<-hang
}
