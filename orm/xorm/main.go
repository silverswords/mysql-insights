package main

import (
	"log"

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

	var user User
	user.Name = "aaa"
	user.Address = "North"
	affected, err := engine.Insert(user)
	if err != nil {
		log.Println(err, "[insert err]")
	}
	log.Println(affected, "[insert affected]")

	result, err := engine.Id(2).Get(&user)
	if err != nil {
		log.Println(err, "[get err]")
	}
	log.Println(result, "[get result]")

	affected, err = engine.Id(1).Delete(user)
	if err != nil {
		log.Println(err, "[delete err]")
	}
	log.Println(affected, "[delete affected]")
}

type User struct {
	Id int64
	// `xorm:"INT(11) NOT NULL AUTO_INCREMENT 'id'"`
	Name    string `xorm:"VARCHAR(64) 'name'"`
	Address string `xorm:"VARCHAR(256) 'address'"`
}
