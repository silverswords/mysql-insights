package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func main() {
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/masterSlaveDB?parseTime=true")
	if err != nil {
		log.Println(err, "err")
	}

	err = engine.Sync2(new(User))
	if err != nil {
		log.Println(err, "[CreateTable err]")
	}

	// user := User{Name: "aaa", Address: "bbb"}
	// start := time.Now()
	// for i := 0; i < 100000; i++ {
	// user := User{Name: "aaa", Address: "bbb"}
	// 	_, err := engine.Insert(user)
	// 	log.Println("[current count]", i)
	// 	log.Println("[current insert time]", time.Now().Sub(start).Seconds())
	// 	if err != nil {
	// 		log.Println(err, "[insert err]")
	// 	}
	// }

	// start := time.Now()
	// result, err := engine.Id(100000).Get(&user)
	// log.Println("[query time]", time.Now().Sub(start).Seconds())
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(result, "[query result]")

	// start := time.Now()
	// affected, err := engine.Update(&User{Name: "ccc", Address: "ddd"})
	// log.Println("[update time]", time.Now().Sub(start).Seconds())
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(affected)

	start := time.Now()
	_, err = engine.Exec("delete from user")
	log.Println("[delete time]", time.Now().Sub(start).Seconds())
	if err != nil {
		log.Println(err)
	}
}

type User struct {
	Id int64
	// `xorm:"INT(11) NOT NULL AUTO_INCREMENT 'id'"`
	Name    string `xorm:"VARCHAR(64) 'name'"`
	Address string `xorm:"VARCHAR(256) 'address'"`
}
